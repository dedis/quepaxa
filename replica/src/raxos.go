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
	"time"
)

const incomingRequestBufferSize = 2000
const numOutgoingThreads = 200
const incomingBufferSize = 1000000
const outgoingBufferSize = int(incomingBufferSize / numOutgoingThreads)

type Instance struct {
	nodeName    int64
	numReplicas int64
	numClients  int64

	//lock sync.Mutex

	replicaAddrList        []string   // array with the IP:port address of every replica
	replicaConnections     []net.Conn // cache of replica connections to all other replicas
	incomingReplicaReaders []*bufio.Reader
	outgoingReplicaWriters []*bufio.Writer

	clientAddrList        []string   // array with the IP:port address of every client
	clientConnections     []net.Conn // cache of client connections to all other clients
	incomingClientReaders []*bufio.Reader
	outgoingClientWriters []*bufio.Writer

	Listener net.Listener // listening to replicas and clients

	rpcTable     map[uint8]*RPCPair
	incomingChan chan *RPCPair

	clientRequestBatchRpc   uint8 // 0
	clientResponseBatchRpc  uint8 // 1
	genericConsensusRpc     uint8 // 2
	messageBlockRpc         uint8 // 3
	messageBlockRequestRpc  uint8 // 4
	clientStatusRequestRpc  uint8 // 5
	clientStatusResponseRpc uint8 // 6
	messageBlockAckRpc      uint8 // 7

	replicatedLog []Slot
	stateMachine  *benchmark.App

	committedIndex int64
	proposedIndex  int64

	proposed []string // assigns the proposed request to the slot

	logFilePath string
	serviceTime int64

	responseSize   int64
	responseString string

	batchSize int64
	batchTime int64

	pipelineLength      int64
	numInflightRequests int64

	outgoingMessageChan chan *OutgoingRPC

	requestsIn   chan *proto.ClientRequestBatch
	messageStore MessageStore
	blockCounter int64

	// from here
	leaderTimeout int64       // in milli seconds
	lastSeenTime  []time.Time // time each replica was last seen

	debugOn bool
}

func New(cfg *configuration.InstanceConfig, name int64, logFilePath string, serviceTime int64, responseSize int64, batchSize int64, batchTime int64, leaderTimeout int64, pipelineLength int64, benchmarkNumber int64, numKeys int64) *Instance {
	in := Instance{
		nodeName:                name,
		numReplicas:             int64(len(cfg.Peers)),
		numClients:              int64(len(cfg.Clients)),
		replicaAddrList:         getReplicaAddressList(cfg), // from here
		replicaConnections:      make([]net.Conn, len(cfg.Peers)),
		incomingReplicaReaders:  make([]*bufio.Reader, len(cfg.Peers)),
		outgoingReplicaWriters:  make([]*bufio.Writer, len(cfg.Peers)),
		clientAddrList:          getClientAddressList(cfg),
		clientConnections:       make([]net.Conn, len(cfg.Clients)),
		incomingClientReaders:   make([]*bufio.Reader, len(cfg.Clients)),
		outgoingClientWriters:   make([]*bufio.Writer, len(cfg.Clients)),
		Listener:                nil,
		rpcTable:                make(map[uint8]*RPCPair),
		incomingChan:            make(chan *RPCPair, incomingBufferSize),
		clientRequestBatchRpc:   0,
		clientResponseBatchRpc:  1,
		genericConsensusRpc:     2,
		messageBlockRpc:         3,
		messageBlockRequestRpc:  4,
		clientStatusRequestRpc:  5,
		clientStatusResponseRpc: 6,
		messageBlockAckRpc:      7,
		//replicatedLog:           nil,
		stateMachine:   benchmark.InitApp(benchmarkNumber, serviceTime, numKeys),
		committedIndex: -1,
		proposedIndex:  -1,
		//proposed:                nil,
		logFilePath:         logFilePath,
		serviceTime:         serviceTime,
		responseSize:        responseSize,
		responseString:      getStringOfSizeN(int(responseSize)),
		batchSize:           batchSize,
		batchTime:           batchTime,
		pipelineLength:      pipelineLength,
		numInflightRequests: 0,
		outgoingMessageChan: make(chan *OutgoingRPC, outgoingBufferSize),
		requestsIn:          make(chan *proto.ClientRequestBatch, incomingRequestBufferSize),
		messageStore:        MessageStore{},
		blockCounter:        0,
		leaderTimeout:       leaderTimeout,
		lastSeenTime:        make([]time.Time, len(cfg.Peers)),
		debugOn:             false,
	}
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

	pid := os.Getpid()
	fmt.Printf("initialized Raxos with process id: %v \n", pid)
	return &in

}
