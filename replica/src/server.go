package raxos

import (
	"raxos/configuration"
	"raxos/proto"
	"time"
)

// server is the main struct for the replica that has a proxy, multiple proposers and a recorder in it

type Server struct {
	ProxyInstance     *Proxy
	ProposerInstances []*Proposer
	RecorderInstance  *Recorder

	proxyToProposerChan  chan ProposeRequest
	proposerToProxyChan  chan ProposeResponse
	recorderToServerChan chan Decision

	lastSeenTimeProposers []time.Time // last seen times of each proposer
}

// ProposeRequest is the message type sent from proxy to proposer

type ProposeRequest struct {
	instance             int64               // slot index
	proposalStr          []string            // fast path client batch ids
	proposalBtch         []proto.ClientBatch // client batches for slow path
	msWait               int                 // number of milliseconds to wait before proposing
	uniqueID             string              // unique id of the proposal
	lastDecidedIndexes   []int               //slots that were previously decided
	lastDecidedDecisions [][]string          // for each lastDecidedIndex, the string array of client batches decided
	lastDecidedUniqueIds []string            // unique id of last decided ids
}

// ProposeResponse is the message type sent from proposer to proxy

type ProposeResponse struct {
	index     int      // log instance
	decisions []string // ids of the client batches
	uniqueId  string
}

type Decision struct {
	indexes   []int      // decided indexes
	uniqueIDs []string   // unique ids of decided indexes
	decisions [][]string // for each index the set of decided client batches
}

// listen to proxy tcp connections, listen to recorder gRPC connections

func (s *Server) NetworkInit() {
	s.ProxyInstance.NetworkInit()    // listen to the proxy port
	s.RecorderInstance.NetworkInit() // listen to gRPC connections
}

// run the main proxy thread which handles all the channels

func (s *Server) Run() {
	go s.ProxyInstance.Run()
}

/*
	start the set of gRPC connections and initiate the set of Proposers
*/
func (s *Server) StartProposers() {
	// create N gRPC connections
	s.setupgRPC()
	// create M number of the Proposers. each have a pointer to the gRPC connections
	s.ProposerInstances = s.createProposers()
	// run each proposer in a separate thread
	s.startProposers()
}

// setup gRPC clients to all recorders and return the connection pointers
func (s *Server) setupgRPC() {
	// todo
}

// create M number of proposers

func (s *Server) createProposers() []*Proposer {
	// todo
	return nil
}

// start M number of proposers

func (s *Server) startProposers() {
	// todo
	// start proposers in separate threads
}

/*
	create a new server instance, inside which there are proxy instance, proposer instances and recorder instance. initialize all fields
*/

func New(cfg *configuration.InstanceConfig, name int64, logFilePath string, batchSize int64, batchTime int64, leaderTimeout int64, pipelineLength int64, benchmark int64, debugOn bool, debugLevel int, leaderMode int) *Server {

	sr := Server{
		ProxyInstance:         nil, //todo
		ProposerInstances:     nil, // todo
		RecorderInstance:      nil, //todo
		proxyToProposerChan:   make(chan ProposeRequest, 10000),
		proposerToProxyChan:   make(chan ProposeResponse, 10000),
		recorderToServerChan:  make(chan Decision, 10000),
		lastSeenTimeProposers: make([]time.Time, len(cfg.Peers)),
	}

	return &sr
}