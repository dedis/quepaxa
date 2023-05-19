package raxos

import (
	"fmt"
	"google.golang.org/grpc"
	"os"
	"raxos/configuration"
	"raxos/proto/client"
	"strconv"
)

// server is the main struct for the replica that has a proxy, multiple proposers and a recorder in it

type Server struct {
	name              int64
	ProxyInstance     *Proxy
	ProposerInstances []*Proposer
	RecorderInstance  *Recorder

	proxyToProposerChan chan ProposeRequest
	proposerToProxyChan chan ProposeResponse
	recorderToProxyChan chan Decision

	proxyToProposerFetchChan chan FetchRequest
	proposerToProxyFetchChan chan FetchResposne

	cfg                         configuration.InstanceConfig // configuration of clients and replicas
	numProposers                int                          // number of proposers == pipeline length
	store                       *ClientBatchStore            // shared client batch store
	serverMode                  int                          // 0 for non-lan-optimized, 1 for lan optimized
	debugOn                     bool
	debugLevel                  int
	leaderTimeout               int64
	leaderMode                  int //0 for fixed leader order, 1 for round robin, static partition,  2 for M.A.B based on commit times, 3 for asynchronous, 4 for last decided proposer
	batchTime                   int64
	epochSize                   int
	proxyToProposerDecisionChan chan Decision
}

// from proxy to proposer

type FetchRequest struct {
	ids []string
}

// from proposer to proxy

type FetchResposne struct {
	batches []client.ClientBatch
}

// ProposeRequest is the message type sent from proxy to proposer

type ProposeRequest struct {
	instance             int64                // slot index
	proposalStr          []string             // fast path client batch ids
	proposalBtch         []client.ClientBatch // client batches for slow path
	isLeader             bool                 // true if waiting time is zero and not in the asyn mode
	lastDecidedIndexes   []int                //slots that were previously decided
	lastDecidedProposers []int32
	lastDecidedDecisions [][]string // for each lastDecidedIndex, the string array of client batches decided
}

// ProposeResponse is the message type sent from proposer to proxy

type ProposeResponse struct {
	index     int      // log instance
	decisions []string // ids of the client batches
	proposer  int32
	s         int
}

type Decision struct {
	indexes   []int      // decided indexes
	decisions [][]string // for each index the set of decided client batches
	proposers []int32
}

// gRPC clients

type peer struct {
	name   int64
	client ConsensusClient
}

// listen to proxy tcp connections, listen to recorder gRPC connections

func (s *Server) NetworkInit() {
	s.ProxyInstance.NetworkInit()    // listen to the proxy port
	s.RecorderInstance.NetworkInit() // listen to gRPC connections
}

// run the main proxy thread which handles all the channels

func (s *Server) Run() {
	s.ProxyInstance.Run()
}

/*
	start the set of gRPC connections and initiate the set of Proposers
*/
func (s *Server) StartProposers() {
	// create M number of the Proposers. each have a pointer to the gRPC connections
	s.createProposers()
}

// setup gRPC clients to all recorders and return the connection pointers

func (s *Server) setupgRPC() []peer {
	peers := make([]peer, 0)
	for _, peeri := range s.cfg.Peers {
		strAddress := peeri.IP + ":" + peeri.RECORDERPORT
		conn, err := grpc.Dial(strAddress, grpc.WithInsecure())
		if err != nil {
			fmt.Fprintf(os.Stderr, "dial: %v", err)
			os.Exit(1)
		}
		intName, _ := strconv.Atoi(peeri.Name)
		newClient := NewConsensusClient(conn)
		peers = append(peers, peer{
			name:   int64(intName),
			client: newClient,
		})
	}
	return peers
}

// create M number of proposers

func (s *Server) createProposers() {
	peers := s.setupgRPC()
	for i := 0; i < s.numProposers+4; i++ { //+4 is for decision sending and fetch requests
		hi := 100000
		newProposer := NewProposer(s.name, int64(i), peers, s.proxyToProposerChan, s.proposerToProxyChan, s.proxyToProposerFetchChan, s.proposerToProxyFetchChan, s.debugOn, s.debugLevel, hi, s.serverMode, s.proxyToProposerDecisionChan)
		s.ProposerInstances = append(s.ProposerInstances, newProposer)
		s.ProposerInstances[len(s.ProposerInstances)-1].runProposer()
	}
}

/*
	create a new server instance, inside which there are proxy instance, proposer instances and recorder instance. initialize all fields
*/

func New(cfg *configuration.InstanceConfig, name int64, logFilePath string, batchSize int64, leaderTimeout int64, pipelineLength int64, debugOn bool, debugLevel int, leaderMode int, serverMode int, batchTime int64, epochSize int, benchmarkMode int, keyLen int, valueLen int, requestPropogationTime int64) *Server {

	sr := Server{
		name:                        name,
		ProxyInstance:               nil,
		ProposerInstances:           nil, // this is initialized in the createProposers method, so no need to create them
		RecorderInstance:            nil,
		proxyToProposerChan:         make(chan ProposeRequest, 10000),
		proposerToProxyChan:         make(chan ProposeResponse, 10000),
		recorderToProxyChan:         make(chan Decision, 10000),
		proxyToProposerFetchChan:    make(chan FetchRequest, 10000),
		proposerToProxyFetchChan:    make(chan FetchResposne, 10000),
		cfg:                         *cfg,
		numProposers:                int(pipelineLength) + 1 + 3,
		store:                       &ClientBatchStore{},
		serverMode:                  serverMode,
		debugOn:                     debugOn,
		debugLevel:                  debugLevel,
		leaderTimeout:               leaderTimeout,
		leaderMode:                  leaderMode,
		batchTime:                   batchTime,
		epochSize:                   epochSize,
		proxyToProposerDecisionChan: make(chan Decision, 10000),
	}

	sr.ProxyInstance = NewProxy(name, *cfg, sr.proxyToProposerChan, sr.proposerToProxyChan, sr.recorderToProxyChan, logFilePath, batchSize, pipelineLength, leaderTimeout, debugOn, debugLevel, &sr, leaderMode, sr.store, serverMode, sr.proxyToProposerFetchChan,
		sr.proposerToProxyFetchChan, sr.batchTime, sr.epochSize, sr.proxyToProposerDecisionChan, benchmarkMode, keyLen, valueLen, requestPropogationTime)
	sr.RecorderInstance = NewRecorder(*cfg, sr.store, sr.recorderToProxyChan, name, debugOn, debugLevel)
	fmt.Printf("started QuePaxa Server")
	return &sr
}
