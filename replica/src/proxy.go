package raxos

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"raxos/common"
	"raxos/configuration"
	"raxos/proto"
	"strconv"
	"sync"
	"time"
)

// slot defines a single instance in the replicated log

type Slot struct {
	// todo
}

/*Proxy saves the state of the proxy which handles client batches, creates replica batches and invokes proposer. Also the proxy executes the SMR and send responses back to client*/

type Proxy struct {
	name        int64 // unique node identifier as defined in the configuration.yml
	numReplicas int
	numClients  int

	clientAddrList        map[int64]string // map with the IP:port address of every client
	incomingClientReaders map[int64]*bufio.Reader
	outgoingClientWriters map[int64]*bufio.Writer

	serverAddress string // proxy address

	buffioWriterMutexes map[int64]*sync.Mutex // to provide mutual exclusion for writes to the same socket connection

	Listener net.Listener // tcp listener for clients

	rpcTable map[uint8]*common.RPCPair

	incomingChan        chan common.RPCPair     // used to collect all the client messages
	outgoingMessageChan chan common.OutgoingRPC // buffer for messages that are written to the wire

	proxyToProposerChan chan ProposeRequest  // proxy to proposer channel
	proposerToProxyChan chan ProposeResponse // proposer to proxy channel

	clientBatchRpc  uint8 // 0
	clientStatusRpc uint8 // 1

	replicatedLog []Slot // the replicated log of the proposer

	exec              bool      // if true the response is sent after execution, if not response is sent after total order
	committedIndex    int64     // last index for which a request was committed and the result was sent to client
	lastProposedIndex int64     // last index proposed
	lastTimeCommitted time.Time // last committed time

	logFilePath string // the path to write the replicated log, used for sanity checks

	batchSize int // maximum replica side batch size
	batchTime int // maximum replica side batch time in micro seconds

	pipelineLength int64 // maximum number of inflight consensus instances

	clientBatchStore ClientBatchStore // message store that stores the client batches

	leaderTimeout int64 // in milliseconds

	debugOn    bool // if turned on, the debug messages will be print on the console
	debugLevel int  // debug level

	serverStarted bool // true if the first status message with operation type 1 received

	server *Server // server instance to call the replica wide functions

	toBeProposed []string // set of client batches that are yet be proposed
}

// instantiate a new proxy

func NewProxy(name int64, cfg configuration.InstanceConfig, proxyToProposerChan chan ProposeRequest, proposerToProxyChan chan ProposeResponse, exec bool, logFilePath string, batchSize int64, batchTime int64, pipelineLength int64, leaderTimeout int64, debugOn bool, debugLevel int, server *Server) *Proxy {

	pr := Proxy{
		name:                  name,
		numReplicas:           len(cfg.Peers),
		numClients:            len(cfg.Clients),
		clientAddrList:        make(map[int64]string),
		incomingClientReaders: make(map[int64]*bufio.Reader),
		outgoingClientWriters: make(map[int64]*bufio.Writer),
		serverAddress:         "",
		buffioWriterMutexes:   make(map[int64]*sync.Mutex),
		Listener:              nil,
		rpcTable:              make(map[uint8]*common.RPCPair),
		incomingChan:          make(chan common.RPCPair),
		outgoingMessageChan:   make(chan common.OutgoingRPC),
		proxyToProposerChan:   proxyToProposerChan,
		proposerToProxyChan:   proposerToProxyChan,
		clientBatchRpc:        0,
		clientStatusRpc:       1,
		replicatedLog:         make([]Slot, 0),
		exec:                  exec,
		committedIndex:        -1,
		lastProposedIndex:     -1,
		logFilePath:           logFilePath,
		batchSize:             int(batchSize),
		batchTime:             int(batchTime),
		pipelineLength:        pipelineLength,
		clientBatchStore:      ClientBatchStore{},
		leaderTimeout:         leaderTimeout,
		debugOn:               debugOn,
		debugLevel:            debugLevel,
		serverStarted:         false,
		server:                server,
		toBeProposed:          make([]string, 0),
	}

	// initialize the clientAddrList

	for i := 0; i < pr.numClients; i++ {
		intName, _ := strconv.Atoi(cfg.Clients[i].Name)
		pr.clientAddrList[int64(intName)] = cfg.Clients[i].IP + ":" + cfg.Clients[i].CLIENTPORT
	}

	// initialize the socketMutexs
	for i := 0; i < len(cfg.Clients); i++ {
		intName, _ := strconv.Atoi(cfg.Clients[i].Name)
		pr.buffioWriterMutexes[int64(intName)] = &sync.Mutex{}
	}

	// serverAddress
	for i := 0; i < len(cfg.Peers); i++ {
		intName, _ := strconv.Atoi(cfg.Peers[i].Name)
		if pr.name == int64(intName) {
			pr.serverAddress = "0.0.0.0:" + cfg.Peers[i].PROXYPORT
		}
	}

	// register rpcs

	pr.RegisterRPC(new(proto.ClientBatch), pr.clientBatchRpc)
	pr.RegisterRPC(new(proto.ClientStatus), pr.clientStatusRpc)

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
			case clientMessage := <-pr.incomingChan:

				//in.lock.Lock()
				pr.debug("Received client  message", 0)
				code := clientMessage.Code
				switch code {
				case pr.clientBatchRpc:
					clientBatch := clientMessage.Obj.(*proto.ClientBatch)
					pr.debug("Client message  "+fmt.Sprintf("%#v", clientBatch), 0)
					pr.ClientBatch(clientBatch)
					break

				case pr.clientStatusRpc:
					clientStatus := clientMessage.Obj.(*proto.ClientStatus)
					pr.debug("Client status  ", 1)
					pr.handleClientStatus(clientStatus)
					break

				}
			case proposerMessage := <-pr.proposerToProxyChan:
				pr.debug("Received proposer message", 0)
				pr.handleProposeResponse(proposerMessage)
			}
		}
	}()
}

/*
	If turned on, print the message to console
*/

func (pr *Proxy) debug(message string, level int) {
	if pr.debugOn && level >= pr.debugLevel {
		fmt.Printf("%s\n", message)
	}
}

// return the pre-agreed, non changing waiting time for the instance by the proposer

func (pr *Proxy) getLeaderWait(instance int) int {
	// todo

	if pr.name == 0 {
		return 0
	} else {
		return 10
	}
}
