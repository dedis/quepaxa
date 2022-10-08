package raxos

import (
	"fmt"
	"google.golang.org/grpc"
	"os"
	"raxos/configuration"
	"raxos/proto"
	"raxos/proto/consensus"
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

	lastSeenTimeProposers []*time.Time // last seen times of each proposer

	peers        []peer                       // set of out going gRPC connections
	cfg          configuration.InstanceConfig // configuration of clients and replicas
	numProposers int                          // number of proposers == pipeline length
	store        *ClientBatchStore            // shared client batch store
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
	uniqueId  string   // unique proposal id
}

type Decision struct {
	indexes   []int      // decided indexes
	uniqueIDs []string   // unique ids of decided indexes
	decisions [][]string // for each index the set of decided client batches
}

// gRPC clients

type peer struct {
	name   int64
	client *consensus.ConsensusClient
}

// AddPeer adds a new peer to the peer list
func (sr *Server) AddPeer(name int64, client *consensus.ConsensusClient) error {

	// add peer to the peer list
	sr.peers = append(sr.peers, peer{
		name:   name,
		client: client,
	})
	return nil
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
	s.ProposerInstances = s.createProposers()
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
		newClient := consensus.NewConsensusClient(conn)
		err = s.AddPeer(int64(intName), &newClient)
		if err != nil {
			conn.Close()
			fmt.Fprintf(os.Stderr, "add peer `%v`: %v", peer.Name, err)
			os.Exit(1)
		}
	}
}

// create M number of proposers

func (s *Server) createProposers() []*Proposer {
	for i := 0; i < s.numProposers; i++ {
		newProposer := NewProposer(s.ProxyInstance.name, int64(i), s.peers, s.proxyToProposerChan, s.proposerToProxyChan, s.lastSeenTimeProposers)
		s.ProposerInstances = append(s.ProposerInstances, newProposer)
		s.ProposerInstances[len(s.ProposerInstances)-1].runProposer()
	}
	return nil
}

/*
	create a new server instance, inside which there are proxy instance, proposer instances and recorder instance. initialize all fields
*/

func New(cfg *configuration.InstanceConfig, name int64, logFilePath string, batchSize int64, leaderTimeout int64, pipelineLength int64, benchmark int64, debugOn bool, debugLevel int, leaderMode int, exec bool) *Server {

	sr := Server{
		ProxyInstance:         nil,
		ProposerInstances:     nil, // this is initialized in the createProposers method, so no need to create them
		RecorderInstance:      nil, //todo add initialization
		proxyToProposerChan:   make(chan ProposeRequest, pipelineLength),
		proposerToProxyChan:   make(chan ProposeResponse, 10000),
		recorderToProxyChan:   make(chan Decision, 10000),
		lastSeenTimeProposers: make([]*time.Time, len(cfg.Peers)),
		peers:                 make([]peer, 0),
		cfg:                   *cfg,
		numProposers:          int(pipelineLength),
		store:                 &ClientBatchStore{},
	}

	sr.ProxyInstance = NewProxy(name, *cfg, sr.proxyToProposerChan, sr.proposerToProxyChan, sr.recorderToProxyChan, exec, logFilePath, batchSize, pipelineLength, leaderTimeout, debugOn, debugLevel, &sr, leaderMode, sr.store)

	return &sr
}
