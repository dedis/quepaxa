package raxos

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"raxos/benchmark"
	"raxos/configuration"
	"raxos/proto"
	_ "raxos/proto"
	"sync"
	"time"
)

const incomingRequestBufferSize = 100000 // size of the buffer that collects incoming client request batches to form blocks
const numOutgoingThreads = 200           // number of wire writers: since the I/O writing is expensive we delegate that task to a thread pool and separate from the critical path
const incomingBufferSize = 1000000       // the size of the buffer which receives all the incoming messages
const outgoingBufferSize = 1000000       // size of the buffer that collects messages to be written to the wire
const hashChannelSize = 10000            // size of the buffer that collects hashes to be sent to the leader

type Instance struct {
	nodeName    int64 // unique node identifier as defined in the configuration.yml
	numReplicas int64 // number of replicas (a replica acts as a proposer and a recorder)
	numClients  int64 // number of clients (this should be known apriori in order to establish tcp connections, since we don't use gRPC)

	//lock sync.Mutex // todo for the moment we don't need this because the shared state is accessed only by the single main thread, but have to verify this

	replicaAddrList []string // array with the IP:port address of every replica
	//replicaConnections     []net.Conn // cache of replica connections to all other replicas
	incomingReplicaReaders []*bufio.Reader
	outgoingReplicaWriters []*bufio.Writer

	clientAddrList []string // array with the IP:port address of every client
	//clientConnections     []net.Conn // cache of client connections to all other clients
	incomingClientReaders []*bufio.Reader
	outgoingClientWriters []*bufio.Writer

	buffioWriterMutexes []sync.Mutex // to provide mutual exclusion for writes to the same socket connection

	Listener net.Listener // tcp listener for replicas and clients

	rpcTable     map[uint8]*RPCPair
	incomingChan chan *RPCPair // used to collect all the incoming messages

	clientRequestBatchRpc   uint8
	clientResponseBatchRpc  uint8
	genericConsensusRpc     uint8
	messageBlockRpc         uint8
	messageBlockRequestRpc  uint8
	clientStatusRequestRpc  uint8
	clientStatusResponseRpc uint8
	messageBlockAckRpc      uint8
	consensusRequestRpc     uint8

	recorderReplicatedLog []Slot         // the replicated log of the recorder
	proposerReplicatedLog []Slot         // the replicated log of the proposer
	stateMachine          *benchmark.App // the application

	committedIndex int64 // last index for which a request was committed and the result was sent to client
	proposedIndex  int64 // last index for which a request was proposed //todo think about the relationship between committed index and the proposed index

	proposed []string // assigns the proposed request to the slot

	logFilePath string // the path to write the replicated log, used for sanity checks
	serviceTime int64  // artificial service time for the no-op app

	responseSize   int64  // fixed response size (might not be useful if the replica doesn't send fixed sized responses)
	responseString string // fixed response string to use if the response size is fixed (might not be used)

	batchSize int64 // maximum server side batch size
	batchTime int64 // maximum replica side batch time in micro seconds

	pipelineLength      int64 // maximum number of inflight consensus instances
	numInflightRequests int64 // current numInflight requests

	outgoingMessageChan chan *OutgoingRPC // buffer for messages that are written to the wire

	requestsIn   chan *proto.ClientRequestBatch // buffer collecting incoming client requests to form blocks
	messageStore MessageStore                   // message store that stores the blocks
	blockCounter int64                          // local sequence number that is used to generate the hash of a block (unique block hash == nodename.blockcounter)

	leaderTimeout int64       // in milliseconds
	lastSeenTime  []time.Time // time each replica was last seen

	debugOn       bool // if turned on, the debug messages will be print on the console
	debugLevel    int  // debug level
	serverStarted bool // true if the first status message with operation type 1 received

	proposeMessage        int64
	spreadEMessage        int64
	spreadCgatherEMessage int64
	gatherCMessage        int64
	decideMessage         int64
	commitMessage         int64

	consensusMessageRecorderDestination int64
	consensusMessageProposerDestination int64
	consensusMessageCommonDestination   int64

	Hi int64 // highest priority for the consensus proposals

	hashProposalsIn chan string // buffer collecting hash values that should be sent to the leader to get proposed
	hashBatchTime   int
	hashBatchSize   int
}

/*

	Instantiate a new Instance object, allocates the buffers
	Initializes the message store

*/

func New(cfg *configuration.InstanceConfig, name int64, logFilePath string, serviceTime int64, responseSize int64, batchSize int64, batchTime int64, leaderTimeout int64, pipelineLength int64, benchmarkNumber int64, numKeys int64, hashBatchSize int64, hashBatchTime int64) *Instance {
	in := Instance{
		nodeName:                name,
		numReplicas:             int64(len(cfg.Peers)),
		numClients:              int64(len(cfg.Clients)),
		replicaAddrList:         GetReplicaAddressList(cfg),
		incomingReplicaReaders:  make([]*bufio.Reader, len(cfg.Peers)),
		outgoingReplicaWriters:  make([]*bufio.Writer, len(cfg.Peers)),
		clientAddrList:          getClientAddressList(cfg),
		incomingClientReaders:   make([]*bufio.Reader, len(cfg.Clients)),
		outgoingClientWriters:   make([]*bufio.Writer, len(cfg.Clients)),
		buffioWriterMutexes:     make([]sync.Mutex, len(cfg.Peers)+len(cfg.Clients)),
		Listener:                nil,
		rpcTable:                make(map[uint8]*RPCPair),
		incomingChan:            make(chan *RPCPair, incomingBufferSize),
		clientRequestBatchRpc:   1,
		clientResponseBatchRpc:  2,
		genericConsensusRpc:     3,
		messageBlockRpc:         4,
		messageBlockRequestRpc:  5,
		clientStatusRequestRpc:  6,
		clientStatusResponseRpc: 7,
		messageBlockAckRpc:      8,
		consensusRequestRpc:     9,
		//replicatedLog:           nil,
		stateMachine:   benchmark.InitApp(benchmarkNumber, serviceTime, numKeys),
		committedIndex: -1,
		proposedIndex:  -1,
		//proposed:                nil,
		logFilePath:                         logFilePath,
		serviceTime:                         serviceTime,
		responseSize:                        responseSize,
		responseString:                      getStringOfSizeN(int(responseSize)),
		batchSize:                           batchSize,
		batchTime:                           batchTime,
		pipelineLength:                      pipelineLength,
		numInflightRequests:                 0,
		outgoingMessageChan:                 make(chan *OutgoingRPC, outgoingBufferSize),
		requestsIn:                          make(chan *proto.ClientRequestBatch, incomingRequestBufferSize),
		messageStore:                        MessageStore{},
		blockCounter:                        0,
		leaderTimeout:                       leaderTimeout,
		lastSeenTime:                        make([]time.Time, len(cfg.Peers)),
		debugOn:                             false,
		debugLevel:                          2, // manually set the debug level
		serverStarted:                       false,
		proposeMessage:                      0,
		spreadEMessage:                      1,
		spreadCgatherEMessage:               2,
		gatherCMessage:                      3,
		decideMessage:                       4,
		commitMessage:                       5,
		consensusMessageRecorderDestination: 0,
		consensusMessageProposerDestination: 1,
		consensusMessageCommonDestination:   2,
		Hi:                                  100000,
		hashProposalsIn:                     make(chan string, hashChannelSize),
		hashBatchSize:                       int(hashBatchSize), // todo this might have to be changed in the WAN
		hashBatchTime:                       int(hashBatchTime),
	}

	for i := 0; i < len(cfg.Peers)+len(cfg.Clients); i++ {
		in.buffioWriterMutexes[i] = sync.Mutex{}
	}

	in.debug("Initialized a new instance", 0)

	rand.Seed(time.Now().UTC().UnixNano())
	in.messageStore.Init()
	/**/
	in.RegisterRPC(new(proto.ClientRequestBatch), in.clientRequestBatchRpc)
	in.RegisterRPC(new(proto.ClientResponseBatch), in.clientResponseBatchRpc)
	in.RegisterRPC(new(proto.GenericConsensus), in.genericConsensusRpc)
	in.RegisterRPC(new(proto.MessageBlock), in.messageBlockRpc)
	in.RegisterRPC(new(proto.MessageBlockRequest), in.messageBlockRequestRpc)
	in.RegisterRPC(new(proto.ClientStatusRequest), in.clientStatusRequestRpc)
	in.RegisterRPC(new(proto.ClientStatusResponse), in.clientStatusResponseRpc)
	in.RegisterRPC(new(proto.MessageBlockAck), in.messageBlockAckRpc)
	in.RegisterRPC(new(proto.ConsensusRequest), in.consensusRequestRpc)

	pid := os.Getpid()
	fmt.Printf("initialized Raxos with process id: %v \n", pid)
	return &in
}
