package raxos

import (
	"fmt"
	"google.golang.org/grpc"
	"os"
	"raxos/configuration"
	"raxos/proto/client"
	"strconv"
	"time"
)

// server is the main struct for the replica that has a proxy, multiple proposers and a recorder in it

type Server struct {
	ProxyInstance     *Proxy
	ProposerInstances []*Proposer
	RecorderInstance  *Recorder

	proxyToProposerChan chan ProposeRequest
	proposerToProxyChan chan ProposeResponse
	recorderToProxyChan chan Decision

	proxyToProposerFetchChan chan FetchRequest
	proposerToProxyFetchChan chan FetchResposne

	lastSeenTimeProposers []*time.Time // last seen times of each proposer

	peers        []peer                       // set of out going gRPC connections
	cfg          configuration.InstanceConfig // configuration of clients and replicas
	numProposers int                          // number of proposers == pipeline length
	store        *ClientBatchStore            // shared client batch store
	serverMode   int
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
	msWait               int                  // number of milliseconds to wait before proposing
	lastDecidedIndexes   []int                //slots that were previously decided
	lastDecidedDecisions [][]string           // for each lastDecidedIndex, the string array of client batches decided
}

// ProposeResponse is the message type sent from proposer to proxy

type ProposeResponse struct {
	index     int      // log instance
	decisions []string // ids of the client batches
}

type Decision struct {
	indexes   []int      // decided indexes
	decisions [][]string // for each index the set of decided client batches
}

// gRPC clients

type peer struct {
	name   int64
	client *ConsensusClient
}

// AddPeer adds a new peer to the peer list
func (sr *Server) AddPeer(name int64, client *ConsensusClient) {
	// add peer to the peer list
	sr.peers = append(sr.peers, peer{
		name:   name,
		client: client,
	})
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
	// create N gRPC connections
	s.setupgRPC()
	// create M number of the Proposers. each have a pointer to the gRPC connections
	s.createProposers()
}

// setup gRPC clients to all recorders and return the connection pointers

func (s *Server) setupgRPC() {
	// add peers
	for _, peer := range s.cfg.Peers {
		strAddress := peer.IP + ":" + peer.RECORDERPORT
		conn, err := grpc.Dial(strAddress, grpc.WithInsecure())
		if err != nil {
			fmt.Fprintf(os.Stderr, "dial: %v", err)
			os.Exit(1)
		}
		intName, _ := strconv.Atoi(peer.Name)
		newClient := NewConsensusClient(conn)
		s.AddPeer(int64(intName), &newClient)
	}
}

// create M number of proposers

func (s *Server) createProposers() {
	for i := 0; i < s.numProposers; i++ {
		newProposer := NewProposer(s.ProxyInstance.name, int64(i), s.peers, s.proxyToProposerChan, s.proposerToProxyChan, s.lastSeenTimeProposers)
		s.ProposerInstances = append(s.ProposerInstances, newProposer)
		s.ProposerInstances[len(s.ProposerInstances)-1].runProposer()
	}
}

/*
	create a new server instance, inside which there are proxy instance, proposer instances and recorder instance. initialize all fields
*/

func New(cfg *configuration.InstanceConfig, name int64, logFilePath string, batchSize int64, leaderTimeout int64, pipelineLength int64, benchmark int64, debugOn bool, debugLevel int, leaderMode int, serverMode int) *Server {

	sr := Server{
		ProxyInstance:            nil,
		ProposerInstances:        nil, // this is initialized in the createProposers method, so no need to create them
		RecorderInstance:         nil,
		proxyToProposerChan:      make(chan ProposeRequest, 10000),
		proposerToProxyChan:      make(chan ProposeResponse, 10000),
		recorderToProxyChan:      make(chan Decision, 10000),
		lastSeenTimeProposers:    make([]*time.Time, len(cfg.Peers)),
		peers:                    make([]peer, 0),
		cfg:                      *cfg,
		numProposers:             int(pipelineLength),
		store:                    &ClientBatchStore{},
		serverMode:               serverMode,
		proxyToProposerFetchChan: make(chan FetchRequest, 10000),
		proposerToProxyFetchChan: make(chan FetchResposne, 10000),
	}

	// allocate the lastSeenTimeProposers

	for i := 0; i < len(sr.lastSeenTimeProposers); i++ {
		sr.lastSeenTimeProposers[i] = &time.Time{}
	}

	sr.ProxyInstance = NewProxy(name, *cfg, sr.proxyToProposerChan, sr.proposerToProxyChan, sr.recorderToProxyChan, logFilePath, batchSize, pipelineLength, leaderTimeout, debugOn, debugLevel, &sr, leaderMode, sr.store, serverMode, sr.proxyToProposerFetchChan, sr.proposerToProxyFetchChan)
	sr.RecorderInstance = NewRecorder(*cfg, sr.store, sr.lastSeenTimeProposers, sr.recorderToProxyChan, name)
	return &sr
}
