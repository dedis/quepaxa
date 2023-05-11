package raxos

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"raxos/benchmark"
	"raxos/common"
	"raxos/configuration"
	"raxos/proto/client"
	"strconv"
	"sync"
	"time"
)

// slot defines a single instance in the replicated log

type Slot struct {
	// slot index is implied by the array position
	proposedBatch []string // client batch ids proposed
	decidedBatch  []string // decided client batch ids
	decided       bool     // true if decided
	committed     bool     // true if committed
	proposer      int32    // whose proposal won
	s             int      // the slot in which proposer decided the entry
}

/*
	proxy saves the state of the proxy and handles client batches, creates replica batches and sends to proposers. Also the proxy executes the SMR and send responses back to client
*/

type Proxy struct {
	name        int64 // unique node identifier as defined in the configuration.yml
	numReplicas int
	numClients  int

	clientAddrList        map[int64]string // map with the IP:port address of every client
	incomingClientReaders map[int64]*bufio.Reader
	outgoingClientWriters map[int64]*bufio.Writer
	buffioWriterMutexes   map[int64]*sync.Mutex // to provide mutual exclusion for writes to the same socket connection

	serverAddress string       // proxy address
	Listener      net.Listener // tcp listener for clients

	rpcTable            map[uint8]*common.RPCPair
	incomingChan        chan common.RPCPair     // used to collect all the client messages
	outgoingMessageChan chan common.OutgoingRPC // buffer for messages that are written to the wire

	proxyToProposerChan      chan ProposeRequest  // proxy to proposer channel
	proposerToProxyChan      chan ProposeResponse // proposer to proxy channel
	recorderToProxyChan      chan Decision        // recorder to proxy channel
	proxyToProposerFetchChan chan FetchRequest
	proposerToProxyFetchChan chan FetchResposne

	clientBatchRpc  uint8 // 1
	clientStatusRpc uint8 // 2

	replicatedLog []Slot // the replicated log of the proxy

	committedIndex    int64     // last index for which a request was committed and the result was sent to client
	lastProposedIndex int64     // last index proposed
	lastTimeCommitted time.Time // last committed time

	logFilePath string // the path to write the replicated log, used for sanity checks

	batchSize      int   // maximum replica side batch size
	batchTime      int64 // micro seconds
	pipelineLength int64 // maximum number of inflight consensus instances

	clientBatchStore *ClientBatchStore // message store that stores the client batches

	leaderTimeout int64 // in microseconds

	debugOn    bool // if turned on, the debug messages will be print on the console
	debugLevel int  // debug level

	serverStarted bool // true if the first status message with operation type 1 received

	server *Server // server instance to call the replica wide functions

	toBeProposed []string // set of client batches that are yet be proposed

	lastDecidedIndexes   []int //slots that were previously decided
	lastDecidedProposers []int32
	lastDecidedDecisions [][]string // for each lastDecidedIndex, the string array of client batches id decided

	leaderMode int // leader change mode
	serverMode int

	instanceTimeouts    []*common.TimerWithCancel
	proposeRequestIndex chan ProposeRequestIndex

	additionalDelay  int // additional delay to add for proposals
	lastTimeProposed time.Time
	epochSize        int

	epochTimes                  []EpochTime
	proxyToProposerDecisionChan chan Decision
	clientBatchTimer            chan ClientBatchTime

	benchmark               *benchmark.Benchmark
	startTime               time.Time
	checkProposerDuplicates bool
	requestPropogationTime  int64

	clientRequestsChan chan client.ClientBatch
}

type EpochTime struct {
	startTime time.Time
	endTime   time.Time
	started   bool
	ended     bool
}

// instantiate a new proxy

func NewProxy(name int64, cfg configuration.InstanceConfig, proxyToProposerChan chan ProposeRequest, proposerToProxyChan chan ProposeResponse, recorderToProxyChan chan Decision, logFilePath string, batchSize int64, pipelineLength int64, leaderTimeout int64, debugOn bool, debugLevel int, server *Server, leaderMode int, store *ClientBatchStore, serverMode int, proxyToProposerFetchChan chan FetchRequest, proposerToProxyFetchChan chan FetchResposne, batchTime int64, epochSize int, proxyToProposerDecisionChan chan Decision, benchmarkMode int, keyLen int, valueLen int, checkProposerDuplicates bool, requestPropogationTime int64) *Proxy {

	pr := Proxy{
		name:                        name,
		numReplicas:                 len(cfg.Peers),
		numClients:                  len(cfg.Clients),
		clientAddrList:              make(map[int64]string),
		incomingClientReaders:       make(map[int64]*bufio.Reader),
		outgoingClientWriters:       make(map[int64]*bufio.Writer),
		buffioWriterMutexes:         make(map[int64]*sync.Mutex),
		serverAddress:               "",
		Listener:                    nil,
		rpcTable:                    make(map[uint8]*common.RPCPair),
		incomingChan:                make(chan common.RPCPair, 1000000),
		outgoingMessageChan:         make(chan common.OutgoingRPC, 1000000),
		proxyToProposerChan:         proxyToProposerChan,
		proposerToProxyChan:         proposerToProxyChan,
		recorderToProxyChan:         recorderToProxyChan,
		proxyToProposerFetchChan:    proxyToProposerFetchChan,
		proposerToProxyFetchChan:    proposerToProxyFetchChan,
		clientBatchRpc:              1,
		clientStatusRpc:             2,
		replicatedLog:               make([]Slot, 0),
		committedIndex:              0,
		lastProposedIndex:           0,
		lastTimeCommitted:           time.Now(),
		logFilePath:                 logFilePath,
		batchSize:                   int(batchSize),
		batchTime:                   batchTime,
		pipelineLength:              pipelineLength,
		clientBatchStore:            store,
		leaderTimeout:               leaderTimeout,
		debugOn:                     debugOn,
		debugLevel:                  debugLevel,
		serverStarted:               false,
		server:                      server,
		toBeProposed:                make([]string, 0),
		lastDecidedIndexes:          make([]int, 0),
		lastDecidedProposers:        make([]int32, 0),
		lastDecidedDecisions:        make([][]string, 0),
		leaderMode:                  leaderMode,
		serverMode:                  serverMode,                               // for the proposer
		instanceTimeouts:            make([]*common.TimerWithCancel, 1000000), // assumes that number of instances do not exceed 1000000, todo increase if not sufficient
		proposeRequestIndex:         make(chan ProposeRequestIndex, 10000),
		additionalDelay:             0,
		lastTimeProposed:            time.Now(),
		epochSize:                   epochSize,
		epochTimes:                  make([]EpochTime, 0),
		proxyToProposerDecisionChan: proxyToProposerDecisionChan,
		clientBatchTimer:            make(chan ClientBatchTime, 100000),
		benchmark:                   benchmark.Init(benchmarkMode, int32(name), keyLen, valueLen),
		checkProposerDuplicates:     checkProposerDuplicates,
		requestPropogationTime:      requestPropogationTime,
		clientRequestsChan:          make(chan client.ClientBatch, 100000),
	}

	// initialize the genenesis

	pr.replicatedLog = append(pr.replicatedLog, Slot{
		proposedBatch: []string{"nil"},
		decidedBatch:  []string{"nil"},
		decided:       true,
		committed:     true,
		proposer:      1,
		s:             0,
	})

	// initialize the clientAddrList

	for i := 0; i < pr.numClients; i++ {
		intName, _ := strconv.Atoi(cfg.Clients[i].Name)
		pr.clientAddrList[int64(intName)] = cfg.Clients[i].IP + ":" + cfg.Clients[i].CLIENTPORT
	}

	// initialize the socketMutexes
	for i := 0; i < len(cfg.Clients); i++ {
		intName, _ := strconv.Atoi(cfg.Clients[i].Name)
		pr.buffioWriterMutexes[int64(intName)] = &sync.Mutex{}
	}

	// serverAddress
	for i := 0; i < len(cfg.Peers); i++ {
		intName, _ := strconv.Atoi(cfg.Peers[i].Name)
		if pr.name == int64(intName) {
			pr.serverAddress = "0.0.0.0:" + cfg.Peers[i].PROXYPORT
			break
		}
	}
	// add a special request with id nil
	pr.clientBatchStore.Add(client.ClientBatch{
		Sender:   -1,
		Messages: nil,
		Id:       "nil",
	})

	// register rpcs

	pr.RegisterRPC(new(client.ClientBatch), pr.clientBatchRpc)
	pr.RegisterRPC(new(client.ClientStatus), pr.clientStatusRpc)

	rand.Seed(time.Now().UTC().UnixNano())

	//pr.debug("initiazlied a new proxy "+fmt.Sprintf("%v", pr.name), -1)

	return &pr
}

// propose request is an internal notification

type ProposeRequestIndex struct {
	index int64
}

// client batch time assigns each incoming client batch with incoming time

type ClientBatchTime struct {
	batch        client.ClientBatch
	incomingTime time.Time
}

/*
	the main loop of the proxy
*/

func (pr *Proxy) Run() {
	go func() {
		for true {

			select {
			case proposeRequest := <-pr.proposeRequestIndex:
				pr.debug("proxy received internal propose request", 1)
				pr.proposeToIndex(proposeRequest.index)
				break
			case proposerMessage := <-pr.proposerToProxyChan:
				pr.debug("proxy received proposer message", -1)
				pr.handleProposeResponse(proposerMessage)
				break
			case recorderMessage := <-pr.recorderToProxyChan:
				//pr.debug("proxy received recorder decide message"+fmt.Sprintf("%v", recorderMessage), -1)
				pr.handleRecorderResponse(recorderMessage)
				break
			case fetchResponse := <-pr.proposerToProxyFetchChan:
				pr.debug("proxy received fetch response", 1)
				pr.handleFetchResponse(fetchResponse)
				break
			case clientRequest := <-pr.clientRequestsChan:
				pr.debug("proxy received a client request that is ready to be replicated", 0)
				pr.handleClientBatch(clientRequest)
				break
			case inpputMessage := <-pr.incomingChan:
				pr.debug("Received client  message", -1)
				code := inpputMessage.Code
				switch code {
				case pr.clientBatchRpc:
					clientBatch := inpputMessage.Obj.(*client.ClientBatch)
					//pr.debug("proxy received client batch  "+fmt.Sprintf("%#v", clientBatch), -1)
					pr.clientBatchTimer <- ClientBatchTime{
						batch:        *clientBatch,
						incomingTime: time.Now(),
					}
					break
				case pr.clientStatusRpc:
					clientStatus := inpputMessage.Obj.(*client.ClientStatus)
					pr.debug("proxy received client status  ", 0)
					pr.handleClientStatus(*clientStatus)
					break
				}
				break
			}

		}
	}()

	// this thread waits for the waiting time to complete for each client request

	go func() {
		for true {
			batch := <-pr.clientBatchTimer
			batchTime := batch.incomingTime
			if time.Now().Sub(batchTime).Milliseconds() >= pr.requestPropogationTime {
				pr.clientRequestsChan <- batch.batch
			} else {
				pr.clientBatchTimer <- batch
			}
		}
	}()

}

/*
	if turned on, print the message to console
*/

func (pr *Proxy) debug(message string, level int) {
	if pr.debugOn && level >= pr.debugLevel {
		fmt.Printf("%s\n", message)
	}
}
