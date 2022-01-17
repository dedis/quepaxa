package cmd

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"os"
	"raxos/configuration"
	"raxos/proto"
	raxos "raxos/replica/src"
	"time"
)

type Client struct {
	clientName  int64 // unique client identifier as defined in the configuration.yml
	numReplicas int64 // number of replicas (a replica acts as a proposer and a recorder)
	numClients  int64 // number of clients (this should be known apriori in order to establish tcp connections, since we don't use gRPC)

	//lock sync.Mutex // todo for the moment we don't need this because the shared state is accessed only by the single main thread, but have to verify this

	replicaAddrList        []string   // array with the IP:port address of every replica
	replicaConnections     []net.Conn // cache of replica connections to all other replicas
	incomingReplicaReaders []*bufio.Reader
	outgoingReplicaWriters []*bufio.Writer

	Listener net.Listener // tcp listener for replicas and clients

	rpcTable     map[uint8]*raxos.RPCPair
	incomingChan chan *raxos.RPCPair // used to collect ClientResponseBatch messages and ClientStatusResponse messages

	clientRequestBatchRpc   uint8 // 0
	clientResponseBatchRpc  uint8 // 1
	clientStatusRequestRpc  uint8 // 2
	clientStatusResponseRpc uint8 // 3

	logFilePath string // the path to write the requests and responses, used for sanity checks

	batchSize int64 // maximum client side batch size
	batchTime int64 // maximum client side batch time

	outgoingMessageChan chan *raxos.OutgoingRPC // buffer for messages that are written to the wire

	requestsIn chan *proto.ClientRequestBatch_SingleClientRequest // buffer collecting incoming client requests to form client batch

	defaultReplica      int64     // id of the default proposer
	replicaTimeout      int64     // in milli seconds: each replica has a unique proposer assigned, upon a timeout, the client changes its default replica
	lastSeenTimeReplica time.Time // time the assigned proposer sent a response

	debugOn bool // if turned on, the debug messages will be print on the console

	requestSize    int64 // size of the request payload in bytes
	testDuration   int64 // test duration in seconds
	warmupDuration int64 // warmup duration in seconds
	arrivalRate    int64 // poisson rate of the request
	benchmark      int64 // type of the workload: 0 for no-op, 1 for kv store and 2 for redis
	numKeys        int64 // maximum number of keys for the key value store

	arrivalTimeChan   chan int64 // channel to which the poisson process adds new request times
	arrivalChan       chan bool  // channel to which the main process adds new request arrivals, to be consumed by the request generation threads
	RequestType       string     // request for sending the client requests, status for sending a status request
	OperationType     int64      // status operation type 1 (bootstrap server), 2: print log
	sentRequests      [][]sentRequestBatch
	receivedResponses []receivedResponseBatch
}

type sentRequestBatch struct {
	batch proto.ClientRequestBatch
	time  time.Time
}

type receivedResponseBatch struct {
	batch proto.ClientResponseBatch
	time  time.Time
}

/*
	Instantiate a new Client object, allocates the buffers
*/

const incomingRequestBufferSize = 2000  // size of the buffer that collects incoming client requests
const numOutgoingThreads = 200          // number of wire writers: since the I/O writing is expensive we delegate that task to a thread pool and separate from the critical path
const numRequestGenerationThreads = 200 // number of  threads that generate client requests upon receiving an arrival indication
const incomingBufferSize = 1000000      // the size of the buffer which receives all the incoming messages (client response bacth messages and client status response message)
const outgoingBufferSize = 1000000      // size of the buffer that collects messages to be written to the wire
const arrivalBufferSize = 1000000       // size of the buffer that collects messages to be written to the wire

func New(name int64, cfg *configuration.InstanceConfig, logFilePath string, batchSize int64, batchTime int64, defaultReplica int64, replicaTimeout int64, requestSize int64, testDuration int64, warmupDuration int64, arrivalRate int64, benchmark int64, numKeys int64, requestType string, operationType int64) *Client {
	cl := Client{
		clientName:              name,
		numReplicas:             int64(len(cfg.Peers)),
		numClients:              int64(len(cfg.Clients)),
		replicaAddrList:         raxos.GetReplicaAddressList(cfg),
		replicaConnections:      make([]net.Conn, len(cfg.Peers)),
		incomingReplicaReaders:  make([]*bufio.Reader, len(cfg.Peers)),
		outgoingReplicaWriters:  make([]*bufio.Writer, len(cfg.Peers)),
		Listener:                nil,
		rpcTable:                make(map[uint8]*raxos.RPCPair),
		incomingChan:            make(chan *raxos.RPCPair, incomingBufferSize),
		clientRequestBatchRpc:   0,
		clientResponseBatchRpc:  1,
		clientStatusRequestRpc:  5,
		clientStatusResponseRpc: 6,
		logFilePath:             logFilePath,
		batchSize:               batchSize,
		batchTime:               batchTime,
		outgoingMessageChan:     make(chan *raxos.OutgoingRPC, outgoingBufferSize),
		requestsIn:              make(chan *proto.ClientRequestBatch_SingleClientRequest, incomingRequestBufferSize),
		defaultReplica:          defaultReplica,
		replicaTimeout:          replicaTimeout,
		lastSeenTimeReplica:     time.Time{},
		debugOn:                 false,
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
		receivedResponses:       make([]receivedResponseBatch, numRequestGenerationThreads),
	}

	rand.Seed(time.Now().UTC().UnixNano())

	/**/
	cl.RegisterRPC(new(proto.ClientRequestBatch), cl.clientRequestBatchRpc)
	cl.RegisterRPC(new(proto.ClientResponseBatch), cl.clientResponseBatchRpc)
	cl.RegisterRPC(new(proto.ClientStatusRequest), cl.clientStatusRequestRpc)
	cl.RegisterRPC(new(proto.ClientStatusResponse), cl.clientStatusResponseRpc)

	pid := os.Getpid()
	fmt.Printf("initialized cLient with process id: %v \n", pid)

	return &cl
}

/*Fill the RPC table by assigning a unique id to each message type*/

func (cl *Client) RegisterRPC(msgObj proto.Serializable, code uint8) {
	cl.rpcTable[code] = &raxos.RPCPair{Code: code, Obj: msgObj}
}

/*
	Each client sends connection requests to all replicas
*/

func (cl *Client) ConnectToReplicas() {
	var b [1]byte
	bs := b[:1]

	//connect to replicas
	for i := int64(0); i < cl.numReplicas; i++ {
		for true {
			conn, err := net.Dial("tcp", cl.replicaAddrList[i])
			if err == nil {
				cl.replicaConnections[i] = conn
				cl.outgoingReplicaWriters[i] = bufio.NewWriter(cl.replicaConnections[i])
				cl.incomingReplicaReaders[i] = bufio.NewReader(cl.replicaConnections[i])
				binary.LittleEndian.PutUint16(bs, uint16(cl.clientName))
				_, err := conn.Write(bs)
				if err != nil {
					panic(err)
				}
				break
			}
		}
	}
	cl.debug("Established all outgoing connections")
}

/*
Listen to all the established tcp connections
*/

func (cl *Client) StartConnectionListners() {
	for i := int64(0); i < cl.numReplicas; i++ {
		go cl.connectionListener(cl.incomingReplicaReaders[i])
	}
}

/*
	listen to a given connection. Upon receiving any message, put it into the central buffer
*/

func (cl *Client) connectionListener(reader *bufio.Reader) {

	var msgType uint8
	var err error = nil

	for true {
		if msgType, err = reader.ReadByte(); err != nil {
			return
		}
		if rpair, present := cl.rpcTable[msgType]; present {
			obj := rpair.Obj.New()
			if err = obj.Unmarshal(reader); err != nil {
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
	This is the main execution thread
	It listens to incoming messages from the incomingChan, and invoke the appropriate handler depending on the message type
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
				cl.debug("Client response batch " + fmt.Sprintf("%#v", clientResponseBatch))
				cl.handleClientResponseBatch(clientResponseBatch)
				break

			case cl.clientStatusResponseRpc:
				clientStatusResponse := replicaMessage.Obj.(*proto.ClientStatusResponse)
				cl.debug("Client Status Response " + fmt.Sprintf("%#v", clientStatusResponse))
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
	msg = oriMsg //
	var w *bufio.Writer

	w = cl.outgoingReplicaWriters[peer]
	err := w.WriteByte(code)
	if err != nil {
		return
	}
	msg.Marshal(w)
	err = w.Flush()
	if err != nil {
		return
	}
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
			}
		}()
	}
}

/*
adds a new out going message to the out going channel
*/

func (cl *Client) sendMessage(peer int64, rpcPair raxos.RPCPair) {
	cl.outgoingMessageChan <- &raxos.OutgoingRPC{
		RpcPair: &rpcPair,
		Peer:    peer,
	}
}
