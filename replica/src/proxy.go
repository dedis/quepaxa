package raxos

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"raxos/common"
	"raxos/configuration"
	"raxos/proto/client"
	"raxos/proto/consensus"
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
}

/*Proxy saves the state of the proxy and handles client batches, creates replica batches and sends to proposers. Also the proxy executes the SMR and send responses back to client*/

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

	clientBatchRpc  uint8 // 0
	clientStatusRpc uint8 // 1
	decideRequest   uint8 // 2
	decideResponse  uint8 // 3

	replicatedLog []Slot // the replicated log of the proxy

	committedIndex    int64     // last index for which a request was committed and the result was sent to client
	lastProposedIndex int64     // last index proposed
	lastTimeCommitted time.Time // last committed time

	logFilePath string // the path to write the replicated log, used for sanity checks

	batchSize      int   // maximum replica side batch size
	pipelineLength int64 // maximum number of inflight consensus instances

	clientBatchStore *ClientBatchStore // message store that stores the client batches

	leaderTimeout int64 // in milliseconds

	debugOn    bool // if turned on, the debug messages will be print on the console
	debugLevel int  // debug level

	serverStarted bool // true if the first status message with operation type 1 received

	server *Server // server instance to call the replica wide functions

	toBeProposed []string // set of client batches that are yet be proposed

	lastDecidedIndexes   []int      //slots that were previously decided
	lastDecidedDecisions [][]string // for each lastDecidedIndex, the string array of client batches decided

	leaderMode int // leader change mode
	serverMode int
}

// instantiate a new proxy

func NewProxy(name int64, cfg configuration.InstanceConfig, proxyToProposerChan chan ProposeRequest, proposerToProxyChan chan ProposeResponse, recorderToProxyChan chan Decision, logFilePath string, batchSize int64, pipelineLength int64, leaderTimeout int64, debugOn bool, debugLevel int, server *Server, leaderMode int, store *ClientBatchStore, serverMode int, proxyToProposerFetchChan chan FetchRequest, proposerToProxyFetchChan chan FetchResposne) *Proxy {

	pr := Proxy{
		name:                     name,
		numReplicas:              len(cfg.Peers),
		numClients:               len(cfg.Clients),
		clientAddrList:           make(map[int64]string),
		incomingClientReaders:    make(map[int64]*bufio.Reader),
		outgoingClientWriters:    make(map[int64]*bufio.Writer),
		buffioWriterMutexes:      make(map[int64]*sync.Mutex),
		serverAddress:            "",
		Listener:                 nil,
		rpcTable:                 make(map[uint8]*common.RPCPair),
		incomingChan:             make(chan common.RPCPair),
		outgoingMessageChan:      make(chan common.OutgoingRPC),
		proxyToProposerChan:      proxyToProposerChan,
		proposerToProxyChan:      proposerToProxyChan,
		recorderToProxyChan:      recorderToProxyChan,
		proxyToProposerFetchChan: proxyToProposerFetchChan,
		proposerToProxyFetchChan: proposerToProxyFetchChan,
		clientBatchRpc:           0,
		clientStatusRpc:          1,
		decideRequest:            2, // not needed
		decideResponse:           3, //not needed
		replicatedLog:            make([]Slot, 0),
		committedIndex:           -1,
		lastProposedIndex:        -1,
		lastTimeCommitted:        time.Now(),
		logFilePath:              logFilePath,
		batchSize:                int(batchSize),
		pipelineLength:           pipelineLength,
		clientBatchStore:         store,
		leaderTimeout:            leaderTimeout,
		debugOn:                  debugOn,
		debugLevel:               debugLevel,
		serverStarted:            false,
		server:                   server,
		toBeProposed:             make([]string, 0),
		lastDecidedIndexes:       make([]int, 0),
		lastDecidedDecisions:     make([][]string, 0),
		leaderMode:               leaderMode,
		serverMode:               serverMode, // for the proposer
	}

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
	pr.RegisterRPC(new(consensus.DecideRequest), pr.decideRequest)
	pr.RegisterRPC(new(consensus.DecideResponse), pr.decideResponse)

	rand.Seed(time.Now().UTC().UnixNano())

	return &pr
}

/*
	the main loop of the proxy
*/

func (pr *Proxy) Run() {
	go func() {
		for true {

			select {
			case inpputMessage := <-pr.incomingChan:

				pr.debug("Received client  message", 0)
				code := inpputMessage.Code
				switch code {
				case pr.clientBatchRpc:
					clientBatch := inpputMessage.Obj.(*client.ClientBatch)
					pr.debug("Client message  "+fmt.Sprintf("%#v", clientBatch), 0)
					pr.handleClientBatch(*clientBatch)
					break

				case pr.clientStatusRpc:
					clientStatus := inpputMessage.Obj.(*client.ClientStatus)
					pr.debug("Client status  ", 1)
					pr.handleClientStatus(*clientStatus)
					break

				}
				break
			case proposerMessage := <-pr.proposerToProxyChan:
				pr.debug("Received proposer message", 0)
				pr.handleProposeResponse(proposerMessage)
				break

			case recorderMessage := <-pr.recorderToProxyChan:
				pr.debug("Received recorder message", 0)
				pr.handleRecorderResponse(recorderMessage)
				break

			case fetchResponse := <-pr.proposerToProxyFetchChan:
				pr.debug("Received fetch response", 0)
				pr.handleFetchResponse(fetchResponse)
				break
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

// return the pre-agreed, non changing waiting time for the instance by the proposer

func (pr *Proxy) getLeaderWait(instance int) int {

	if pr.leaderMode == 0 {
		// fixed order
		if pr.name == 0 {
			return 0
		} else {
			return int(pr.name * pr.leaderTimeout)
		}
	} else if pr.leaderMode == 1 {
		// todo
		// static MAB
		return 0
	} else if pr.leaderMode == 2 {
		// todo
		// dynamic MAB
		return 0
	}
	return 10
}
