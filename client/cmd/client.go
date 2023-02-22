package cmd

import (
	"bufio"
	"fmt"
	"os"
	"raxos/common"
	"raxos/configuration"
	"raxos/proto/client"
	"strconv"
	"sync"
	"time"
)

/*
	This file defines the client struct and the new method that is invoked when creating a new client by the main

*/

type Client struct {
	name        int64 // unique client identifier as defined in the configuration.yml
	numReplicas int64 // number of replicas
	numClients  int64 // number of clients (for the moment this is not required, just leaving for future use or to remove)

	replicaAddrList        map[int64]string        // a map with the IP:port of every proxy
	incomingReplicaReaders map[int64]*bufio.Reader // socket readers for each replica
	outgoingReplicaWriters map[int64]*bufio.Writer // socket writer for each replica
	socketMutexs           map[int64]*sync.Mutex   // for mutual exclusion for each buffio.writer outgoingReplicaWriters

	rpcTable     map[uint8]*common.RPCPair // map each RPC type (message type) to its unique number
	incomingChan chan *common.RPCPair      // used to collect ClientBatch messages and ClientStatus messages (basically all the incoming messages)

	clientBatchRpc  uint8 // 1
	clientStatusRpc uint8 // 2

	logFilePath string // the path to write the requests and responses, used for sanity checks

	batchSize int64 // maximum client side batch size
	batchTime int64

	outgoingMessageChan chan *common.OutgoingRPC // buffer for messages that are written to the wire

	debugOn    bool // if turned on, the debug messages will be printed on the console
	debugLevel int  // debug level

	testDuration int64 // test duration in seconds
	arrivalRate  int64 // poisson rate of the request (requests per second)

	arrivalTimeChan   chan int64           // channel to which the poisson process adds new request arrival times in nanoseconds w.r.t test start time
	arrivalChan       chan bool            // channel to which the main scheduler adds new request arrivals, to be consumed by the request generation threads
	RequestType       string               // [request] for sending a stream of client requests, [status] for sending a status request
	OperationType     int64                // status operation type 1 (bootstrap server), 2: print log
	sentRequests      [][]sentRequestBatch // generator i updates sentRequests[i] :this is to avoid concurrent access to the same array
	receivedResponses sync.Map             // set of received client response batches from replicas
	startTime         time.Time            // test start time

	clientListenAddress string // ip:port of the listening port

	keyLen         int // length of keys
	valLen         int // length of values
	slowdown       string
	useFixedLeader bool  // if true send only to a fixed leader
	fixedLeader    int32 //
}

/*
	sentRequestBatch contains a batch that was written to wire, and the time it was written
*/

type sentRequestBatch struct {
	batch client.ClientBatch
	time  time.Time
}

/*
	received response batch contains the batch that was received from the replicas
*/

type receivedResponseBatch struct {
	batch client.ClientBatch
	time  time.Time
}

const statusTimeout = 10              // time to wait for a status request in seconds
const numOutgoingThreads = 200        // number of wire writers: since the I/O writing is expensive we delegate that task to a thread pool and separate from the critical path
const numRequestGenerationThreads = 4 // number of  threads that generate client requests upon receiving an arrival indication
const incomingBufferSize = 1000000    // the size of the buffer which receives all the incoming messages (client response batch messages and client status response message)
const outgoingBufferSize = 1000000    // size of the buffer that collects messages to be written to the wire
const arrivalBufferSize = 1000000     // size of the buffer that collects new request arrivals

/*
	Instantiate a new Client instance, allocate the buffers
*/

func New(name int64, cfg *configuration.InstanceConfig, logFilePath string, batchSize int64, batchTime int64, testDuration int64, arrivalRate int64, requestType string, operationType int64, debugLevel int, debugOn bool, keyLen int, valLen int, slowdown string, useFixedLeader bool, fixedLeader int64) *Client {
	cl := Client{
		name:                   name,
		numReplicas:            int64(len(cfg.Peers)),
		numClients:             int64(len(cfg.Clients)),
		replicaAddrList:        make(map[int64]string),
		incomingReplicaReaders: make(map[int64]*bufio.Reader),
		outgoingReplicaWriters: make(map[int64]*bufio.Writer),
		socketMutexs:           make(map[int64]*sync.Mutex),
		rpcTable:               make(map[uint8]*common.RPCPair),
		incomingChan:           make(chan *common.RPCPair, incomingBufferSize),
		clientBatchRpc:         1,
		clientStatusRpc:        2,
		logFilePath:            logFilePath,
		batchSize:              batchSize,
		batchTime:              batchTime,
		outgoingMessageChan:    make(chan *common.OutgoingRPC, outgoingBufferSize),
		debugOn:                debugOn,
		debugLevel:             debugLevel,
		testDuration:           testDuration,
		arrivalRate:            arrivalRate,
		arrivalTimeChan:        make(chan int64, arrivalBufferSize),
		arrivalChan:            make(chan bool, arrivalBufferSize),
		RequestType:            requestType,
		OperationType:          operationType,
		sentRequests:           make([][]sentRequestBatch, numRequestGenerationThreads),
		receivedResponses:      sync.Map{},
		startTime:              time.Time{},
		clientListenAddress:    "",
		keyLen:                 keyLen,
		valLen:                 valLen,
		slowdown:               slowdown,
		useFixedLeader:         useFixedLeader,
		fixedLeader:            int32(fixedLeader),
	}

	// client listen address

	for i := 0; i < len(cfg.Clients); i++ {
		if cfg.Clients[i].Name == strconv.Itoa(int(name)) {
			cl.clientListenAddress = "0.0.0.0:" + cfg.Clients[i].CLIENTPORT
			break
		}
	}

	// initialize the replicaAddrList
	for i := int64(0); i < cl.numReplicas; i++ {
		intName, _ := strconv.Atoi(cfg.Peers[i].Name)
		cl.replicaAddrList[int64(intName)] = cfg.Peers[i].IP + ":" + cfg.Peers[i].PROXYPORT
	}

	// initialize the socketMutexs
	for i := int64(0); i < cl.numReplicas; i++ {
		intName, _ := strconv.Atoi(cfg.Peers[i].Name)
		cl.socketMutexs[int64(intName)] = &sync.Mutex{}
	}

	// initialize sentRequests

	for i := 0; i < numRequestGenerationThreads; i++ {
		cl.sentRequests[i] = make([]sentRequestBatch, 0)
	}

	cl.RegisterRPC(new(client.ClientBatch), cl.clientBatchRpc)
	cl.RegisterRPC(new(client.ClientStatus), cl.clientStatusRpc)

	pid := os.Getpid()
	fmt.Printf("initialized client %v with process id: %v \n", name, pid)
	return &cl
}

/*
	if turned on, prints the message to console
*/

func (cl *Client) debug(message string, level int) {
	if cl.debugOn && level >= cl.debugLevel {
		fmt.Printf("%s\n", message)
	}
}
