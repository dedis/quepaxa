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

	proxyToProposerChan ProposeRequest
	proposerToProxyChan ProposeResponse

	lastSeenTimeProposers []time.Time // last seen times of each proposer
}

// ProposeRequest is the message type sent from proxy to proposer

type ProposeRequest struct {
	instance int64 // slot index
	proposalStr []string // fast path client batch ids
	proposalBtch []*proto.ClientBatch // client batches for slow path
	msWait int // number of milliseconds to wait before proposing 
	uniqueID string // unique id of the proposal
	lastDecidedIndexes []int //slots that were previously decided
	lastDecidedDecisions [][]string // for each lastDecidedIndex, the string array of client batches decided 
	lastDecidedUniqueIds []string // unique id of last decided ids
}

// ProposeResponse is the message type sent from proposer to proxy

type ProposeResponse struct {
	//todo start from here
}

// listen to proxy tcp connections, listen to recorder gRPC connections

func (s Server) NetworkInit() {
	s.ProxyInstance.NetworkInit()
	s.RecorderInstance.NetworkInit()
}

// run the main proxy thread which handles all the channels

func (s Server) Run() {
	go s.ProxyInstance.Run()
}

// start the set of gRPC connections and initiate the set of Proposers
func (s Server) StartgRPC(){
	
}

/*
	create a new server instance, inside which there are proxy instance, proposer instances and recorder instance. initialize all fields
*/

func New(cfg *configuration.InstanceConfig, name int64, logFilePath string, batchSize int64, batchTime int64, leaderTimeout int64, pipelineLength int64, benchmark int64, debugOn bool, debugLevel int) *Server {

	sr := Server{}

	return &sr
}
