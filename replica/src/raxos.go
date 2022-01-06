package raxos

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"raxos/benchmark"
	"raxos/configuration"
	"sync"
	"time"
)

func getRealSizeOf(v interface{}) (int, error) {
	b := new(bytes.Buffer)
	if err := gob.NewEncoder(b).Encode(v); err != nil {
		return 0, err
	}
	return b.Len(), nil
}

func getStringOfSizeN(length int) string {
	str := "a"
	size, _ := getRealSizeOf(str)
	for size < length {
		str = str + "a"
		size, _ = getRealSizeOf(str)
	}
	return str
}

type Instance struct {
	nodeName    int64
	numReplicas int64
	numClients  int64

	lock sync.Mutex

	replicaAddrList        []string   // array with the IP:port address of every replica
	replicaConnections     []net.Conn // cache of replica connections to all other replicas
	incomingReplicaReaders []*bufio.Reader
	outgoingReplicaWriters []*bufio.Writer

	clientAddrList        []string   // array with the IP:port address of every client
	clientConnections     []net.Conn // cache of client connections to all other clients
	incomingClientReaders []*bufio.Reader
	outgoingClientWriters []*bufio.Writer

	Listener net.Listener // listening to replicas and clients

	connectedToReplicas bool
	connectedToClients  bool
	startedServer       bool
	startedHeartBeats   bool
	startedBatch        bool

	rpcTable map[uint8]*RPCPair

	instances        []instance
	nextFreeInstance int64

	logFilePath string
	serviceTime int64

	responseSize   int64
	responseString string

	batchSize   int64
	batchTime   int64
	requestsIn  chan request
	requestsOut chan bool

	replicaChan chan replicaMessage

	clientOutChan chan int

	prepareRpc uint8
	ackRpc     uint8
	proposeRpc uint8
	acceptRpc  uint8
	learnRpc   uint8
	emptyRpc   uint8

	expected  int64
	index     int64
	est_index []int64

	keepAliveTimeout int64 // in milli seconds
	lastSeenTime     []time.Time

	outgoingMessageChannels []chan RPCPair

	debugOn bool
	app     *benchmark.App
}

func (in *Instance) debug(message string) {
	fmt.Printf("%s\n", message)

}

func New(cfg *configuration.InstanceConfig, name int64, logFilePath string, serviceTime int64, responseSize int64, batchSize int64, batchTime int64, leaderTimeout int64, pipelineLength int64, benchmark int64, numKeys int64) {

}
