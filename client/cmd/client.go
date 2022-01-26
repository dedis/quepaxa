package cmd

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"raxos/configuration"
	"raxos/proto"
	raxos "raxos/replica/src"
	"strconv"
	"sync"
	"time"
)

/*
	This file defines the client struct and the new method that is invoked when creating a new client by the main

*/

type Client struct {
	clientName  int64 // unique client identifier as defined in the configuration.yml
	numReplicas int64 // number of replicas (a replica acts as a proposer and a recorder)
	numClients  int64 // number of clients (for the moment this is not required, just leaving for future use or to remove)

	//lock sync.Mutex // todo for the moment we don't need this because the shared state is accessed only by the single main thread, but have to verify this

	replicaAddrList []string // array with the IP:port address of every replica
	//replicaConnections     []net.Conn      // cache of replica connections to all others replicas, for the moment we don't need this
	incomingReplicaReaders []*bufio.Reader // socket readers for each replica
	outgoingReplicaWriters []*bufio.Writer // socket writer for each replica
	socketMutexs           []sync.Mutex    // for mutual exclusion for each buffio.writer outgoingReplicaWriters

	rpcTable     map[uint8]*raxos.RPCPair // map each RPC type (message type) to its unique number
	incomingChan chan *raxos.RPCPair      // used to collect ClientResponseBatch messages and ClientStatusResponse messages (basically all the incoming messages)

	clientRequestBatchRpc   uint8 // 0
	clientResponseBatchRpc  uint8 // 1
	clientStatusRequestRpc  uint8 // 5
	clientStatusResponseRpc uint8 // 6

	/*
		note that the client doesn't receive any consensus message, or block message
	*/

	logFilePath string // the path to write the requests and responses, used for sanity checks

	batchSize int64 // maximum client side batch size
	batchTime int64 // maximum client side batch time in micro seconds

	outgoingMessageChan chan *raxos.OutgoingRPC // buffer for messages that are written to the wire

	defaultReplica      int64     // id of the default proposer to which the client sends the messages
	replicaTimeout      int64     // in seconds: each replica has a unique proposer assigned, upon a timeout, the client changes its default replica
	lastSeenTimeReplica time.Time // time the assigned proposer last sent a response

	debugOn bool // if turned on, the debug messages will be print on the console

	requestSize    int64 // size of the request payload in bytes (applicable only for the no-op application / testing purposes)
	testDuration   int64 // test duration in seconds
	warmupDuration int64 // warmup duration in seconds (at the moment not used)
	arrivalRate    int64 // poisson rate of the request (requests per second)
	benchmark      int64 // type of the workload: 0 for no-op, 1 for kv store and 2 for redis
	numKeys        int64 // maximum number of keys for the key value store

	arrivalTimeChan     chan int64              // channel to which the poisson process adds new request arrival times in nanoseconds w.r.t test start time
	arrivalChan         chan bool               // channel to which the main scheduler adds new request arrivals, to be consumed by the request generation threads
	RequestType         string                  // [request] for sending a stream of client requests, [status] for sending a status request
	OperationType       int64                   // status operation type 1 (bootstrap server), 2: print log
	sentRequests        [][]sentRequestBatch    // generator i updates sentRequests[i] :this is to avoid concurrent access to the same array
	receivedResponses   []receivedResponseBatch // set of received client response batches from replicas
	startTime           time.Time               // test start time
	clientListenAddress string                  // TCP addresses to which the client listens to new incoming TCP connections
}

/*
	sentRequestBatch contains a batch that was written to wire, and the time it was written
*/

type sentRequestBatch struct {
	batch proto.ClientRequestBatch
	time  time.Time
}

/*
	received response batch contains the batch that was received from the replicas
*/

type receivedResponseBatch struct {
	batch proto.ClientResponseBatch
	time  time.Time
}

const statusTimeout = 10              // time to wait for a status request in seconds, todo might have to increase depending on the test duration and throughput
const numOutgoingThreads = 200        // number of wire writers: since the I/O writing is expensive we delegate that task to a thread pool and separate from the critical path
const numRequestGenerationThreads = 4 // number of  threads that generate client requests upon receiving an arrival indication todo try different values for this: lower values result in big batches
const incomingBufferSize = 1000000    // the size of the buffer which receives all the incoming messages (client response batch messages and client status response message)
const outgoingBufferSize = 1000000    // size of the buffer that collects messages to be written to the wire
const arrivalBufferSize = 1000000     // size of the buffer that collects new request arrivals

/*
	Instantiate a new Client instance, allocate the buffers
*/

func New(name int64, cfg *configuration.InstanceConfig, logFilePath string, batchSize int64, batchTime int64, defaultReplica int64, replicaTimeout int64, requestSize int64, testDuration int64, warmupDuration int64, arrivalRate int64, benchmark int64, numKeys int64, requestType string, operationType int64) *Client {
	cl := Client{
		clientName:      name,
		numReplicas:     int64(len(cfg.Peers)),
		numClients:      int64(len(cfg.Clients)),
		replicaAddrList: raxos.GetReplicaAddressList(cfg),
		//replicaConnections:      make([]net.Conn, len(cfg.Peers)),
		incomingReplicaReaders:  make([]*bufio.Reader, len(cfg.Peers)),
		outgoingReplicaWriters:  make([]*bufio.Writer, len(cfg.Peers)),
		socketMutexs:            make([]sync.Mutex, len(cfg.Peers)),
		rpcTable:                make(map[uint8]*raxos.RPCPair),
		incomingChan:            make(chan *raxos.RPCPair, incomingBufferSize),
		clientRequestBatchRpc:   1,
		clientResponseBatchRpc:  2,
		clientStatusRequestRpc:  6,
		clientStatusResponseRpc: 7,
		logFilePath:             logFilePath,
		batchSize:               batchSize,
		batchTime:               batchTime,
		outgoingMessageChan:     make(chan *raxos.OutgoingRPC, outgoingBufferSize),
		defaultReplica:          defaultReplica,
		replicaTimeout:          replicaTimeout,
		lastSeenTimeReplica:     time.Now(),
		debugOn:                 false, // manually set this if debugging needs to be turned on
		requestSize:             requestSize,
		testDuration:            testDuration,
		warmupDuration:          warmupDuration,
		arrivalRate:             arrivalRate,
		benchmark:               benchmark,
		numKeys:                 numKeys,
		arrivalTimeChan:         make(chan int64, arrivalBufferSize),
		arrivalChan:             make(chan bool, arrivalBufferSize),
		RequestType:             requestType,
		OperationType:           operationType,
		sentRequests:            make([][]sentRequestBatch, numRequestGenerationThreads),
		startTime:               time.Now(),
		clientListenAddress:     cfg.Clients[int(name)-len(cfg.Peers)].Address,
	}

	cl.debug("Created a new client instance")

	for i := 0; i < len(cfg.Peers); i++ {
		cl.socketMutexs[i] = sync.Mutex{}
	}

	rand.Seed(time.Now().UTC().UnixNano())

	/*
		Register the rpcs
	*/
	cl.RegisterRPC(new(proto.ClientRequestBatch), cl.clientRequestBatchRpc)
	cl.RegisterRPC(new(proto.ClientResponseBatch), cl.clientResponseBatchRpc)
	cl.RegisterRPC(new(proto.ClientStatusRequest), cl.clientStatusRequestRpc)
	cl.RegisterRPC(new(proto.ClientStatusResponse), cl.clientStatusResponseRpc)

	cl.debug("Registered RPCs in the table")

	pid := os.Getpid()
	fmt.Printf("initialized client %v with process id: %v \n", name, pid)

	return &cl
}

/*
	Fill the RPC table by assigning a unique id to each message type
*/

func (cl *Client) RegisterRPC(msgObj proto.Serializable, code uint8) {
	cl.rpcTable[code] = &raxos.RPCPair{Code: code, Obj: msgObj}
}

/*
	Each client sends connection requests to all replicas
*/

func (cl *Client) ConnectToReplicas() {

	cl.debug("Connecting to " + strconv.Itoa(int(cl.numReplicas)) + " replicas")

	var b [4]byte
	bs := b[:4]

	//connect to replicas
	for i := int64(0); i < cl.numReplicas; i++ {
		for true {
			conn, err := net.Dial("tcp", cl.replicaAddrList[i])
			if err == nil {
				cl.outgoingReplicaWriters[i] = bufio.NewWriter(conn)
				binary.LittleEndian.PutUint16(bs, uint16(cl.clientName))
				_, err := conn.Write(bs)
				if err != nil {
					cl.debug("Error while connecting to replica" + strconv.Itoa(int(i)))
					panic(err)
				}
				cl.debug("Established outgoing connection to " + strconv.Itoa(int(i)))
				break
			}
		}
	}
	cl.debug("Established all outgoing connections to replicas")
}

/*
	Listen on the server port for new connections
	Whenever a replica receives a new client connection, it dials the client
*/

func (cl *Client) WaitForConnections() {

	var b [4]byte
	bs := b[:4]
	Listener, _ := net.Listen("tcp", cl.clientListenAddress)
	cl.debug("Listening to incoming connections on " + cl.clientListenAddress)

	for true {
		conn, err := Listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			panic(err)
		}
		if _, err := io.ReadFull(conn, bs); err != nil {
			fmt.Println("Connection establish error:", err)
			panic(err)
		}
		id := int32(binary.LittleEndian.Uint16(bs))
		cl.debug("Received incoming tcp connection from " + strconv.Itoa(int(id)))

		cl.incomingReplicaReaders[id] = bufio.NewReader(conn)
		go cl.connectionListener(cl.incomingReplicaReaders[id], id)
		cl.debug("Started listening to " + strconv.Itoa(int(id)))

	}
}

/*
	listen to a given connection. Upon receiving any message, put it into the central buffer
*/

func (cl *Client) connectionListener(reader *bufio.Reader, id int32) {

	var msgType uint8
	var err error = nil

	for true {
		if msgType, err = reader.ReadByte(); err != nil {
			cl.debug("Error while reading message code: connection broken from " + strconv.Itoa(int(id)))
			return
		}
		if rpair, present := cl.rpcTable[msgType]; present {
			obj := rpair.Obj.New()
			if err = obj.Unmarshal(reader); err != nil {
				cl.debug("Error while unmarshalling")
				return
			}
			cl.incomingChan <- &raxos.RPCPair{
				Code: msgType,
				Obj:  obj,
			}
		} else {
			cl.debug("Error: received unknown message type")
		}
	}
}

/*
	If turned on, prints the message to console
*/

func (cl *Client) debug(message string) {
	if cl.debugOn {
		fmt.Printf("%s\n", message)
	}
}

/*
	This is an execution thread that listens to all the incoming messages
	It listens to incoming messages from the incomingChan, and invokes the appropriate handler depending on the message type
*/

func (cl *Client) Run() {
	go func() {
		for true {

			cl.debug("Checking channel\n")
			replicaMessage := <-cl.incomingChan
			//in.lock.Lock()
			cl.debug("Received replica message")
			code := replicaMessage.Code
			switch code {

			case cl.clientResponseBatchRpc:
				clientResponseBatch := replicaMessage.Obj.(*proto.ClientResponseBatch)
				cl.debug("Client response batch from " + fmt.Sprintf("%#v", clientResponseBatch.Sender))
				cl.handleClientResponseBatch(clientResponseBatch)
				break

			case cl.clientStatusResponseRpc:
				clientStatusResponse := replicaMessage.Obj.(*proto.ClientStatusResponse)
				cl.debug("Client Status Response from" + fmt.Sprintf("%#v", clientStatusResponse.Sender))
				cl.handleClientStatusResponse(clientStatusResponse)
				break
			}
			//in.lock.Unlock()
		}
	}()
}

/*
	Write a message to the wire, first the message type is written and then the actual message
*/

func (cl *Client) internalSendMessage(peer int64, rpcPair *raxos.RPCPair) {
	code := rpcPair.Code
	oriMsg := rpcPair.Obj
	var msg proto.Serializable
	msg = oriMsg // unlike in the replica where we generate a new message, to avoid unsafe concurrent proto operations, we use a single message object, because the client doesn't broadcast (even when it does it create a new message)
	var w *bufio.Writer

	w = cl.outgoingReplicaWriters[peer]
	cl.socketMutexs[peer].Lock()
	err := w.WriteByte(code)
	if err != nil {
		cl.debug("Error writing message code byte:" + err.Error())
		return
	}
	err = msg.Marshal(w)
	if err != nil {
		cl.debug("Error while marshalling:" + err.Error())
		return
	}
	err = w.Flush()
	if err != nil {
		cl.debug("Error flushing:" + err.Error())
		return
	}
	cl.socketMutexs[peer].Unlock()
	cl.debug("Internal sent message to " + strconv.Itoa(int(peer)))
}

/*
	A set of threads that manages outgoing messages: write the message to the OS buffers
*/

func (cl *Client) StartOutgoingLinks() {
	for i := 0; i < numOutgoingThreads; i++ {
		go func() {
			for true {
				outgoingMessage := <-cl.outgoingMessageChan
				cl.internalSendMessage(outgoingMessage.Peer, outgoingMessage.RpcPair)
				cl.debug("Invoked internal sent to replica " + strconv.Itoa(int(outgoingMessage.Peer)))
			}
		}()
	}
}

/*
	adds a new out-going message to the out going channel
*/

func (cl *Client) sendMessage(peer int64, rpcPair raxos.RPCPair) {
	cl.outgoingMessageChan <- &raxos.OutgoingRPC{
		RpcPair: &rpcPair,
		Peer:    peer,
	}
	cl.debug("Added RPC pair to outgoing channel to peer " + strconv.Itoa(int(peer)))
}
