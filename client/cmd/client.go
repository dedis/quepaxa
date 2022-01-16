package cmd

import (
	"bufio"
	"net"
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

	arrivalChan       chan bool // channel to which the poisson process adds new request times
	requestType       string    // request for sending the client requests, status for sending a status request
	operationType     int64     // status operation type 1 (bootstrap server), 2: print log
	sentRequests      []proto.ClientRequestBatch
	receivedResponses []proto.ClientResponseBatch
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
		arrivalChan:             make(chan bool, arrivalBufferSize),
		requestType:             requestType,
		operationType:           operationType,
	}

	return &cl
}
