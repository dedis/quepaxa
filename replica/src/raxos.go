package raxos

import (
	"bufio"
	"net"
	"raxos/benchmark"
	"raxos/configuration"
	"raxos/internal"
	_ "raxos/proto"
	"time"
)

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

	clientRequestRpc       uint8 // 0
	clientResponseRpc      uint8 // 1
	genericConsensusRpc    uint8 // 2
	messageBlockReplyRpc   uint8 // 3
	messageBlockRequestRpc uint8 // 4

	instances    []internal.Slot
	stateMachine *benchmark.App

	committedIndex int64
	proposedIndex  int64

	pendingRepliesMap map[uint64]int64 // assigns request identifier cm to the client cl that is waiting for the reply
	proposed          []int64          // assigns the proposed request to the slot

	logFilePath string
	serviceTime int64

	responseSize   int64
	responseString string

	batchSize int64
	batchTime int64

	outgoingMessageChan chan *OutgoingRPC

	// from here

	requestsIn  chan request
	requestsOut chan bool

	clientOutChan chan int

	expected  int64
	index     int64
	est_index []int64

	keepAliveTimeout int64 // in milli seconds
	lastSeenTime     []time.Time

	outgoingMessageChannels []chan RPCPair

	debugOn bool
}

func New(cfg *configuration.InstanceConfig, name int64, logFilePath string, serviceTime int64, responseSize int64, batchSize int64, batchTime int64, leaderTimeout int64, pipelineLength int64, benchmark int64, numKeys int64) {

}
