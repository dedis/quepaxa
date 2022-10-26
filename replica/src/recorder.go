package raxos

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
	"os"
	"raxos/configuration"
	"strconv"
	"sync"
	"time"
)

type Recorder struct {
	address               string       // address to listen for gRPC connections
	listener              net.Listener // socket for gRPC connections
	server                *grpc.Server // gRPC server
	connection            *GRPCConnection
	clientBatches         *ClientBatchStore
	lastSeenTimeProposers []*time.Time // last seen times of each proposer
	recorderToProxyChan   chan Decision
	name                  int64
	slots                 []RecorderSlot // recorder side replicated log
	cfg                   configuration.InstanceConfig
	instanceCreationMutex *sync.Mutex
}

type Value struct {
	priority    int64
	proposer_id int64
	thread_id   int64
	ids         []string
}

// Recorder side slot

type RecorderSlot struct {
	Mutex *sync.Mutex
	S     int
	F     Value
	A     Value
	M     Value
}

// instantiate a new Recorder

func NewRecorder(cfg configuration.InstanceConfig, clientBatches *ClientBatchStore, lastSeenTimeProposers []*time.Time, recorderToProxyChan chan Decision, name int64) *Recorder {

	re := Recorder{
		address:               "",
		listener:              nil,
		server:                nil,
		connection:            nil,
		clientBatches:         clientBatches,
		lastSeenTimeProposers: lastSeenTimeProposers,
		recorderToProxyChan:   recorderToProxyChan,
		name:                  name,
		slots:                 make([]RecorderSlot, 0),
		cfg:                   cfg,
		instanceCreationMutex: &sync.Mutex{},
	}

	// initialize the address

	// serverAddress
	for i := 0; i < len(cfg.Peers); i++ {
		intName, _ := strconv.Atoi(cfg.Peers[i].Name)
		if re.name == int64(intName) {
			re.address = "0.0.0.0:" + cfg.Peers[i].RECORDERPORT
			break
		}
	}

	return &re
}

// start listening to gRPC connection

func (r *Recorder) NetworkInit() {
	r.server = grpc.NewServer()
	r.connection = &GRPCConnection{Recorder: r}
	RegisterConsensusServer(r.server, r.connection)
	RegisterFetchServer(r.server, r.connection)

	// start listener
	listener, err := net.Listen("tcp", r.address)
	r.listener = listener
	if err != nil {
		fmt.Fprintf(os.Stderr, "listen: %v", err)
		os.Exit(1)
	}
	go r.server.Serve(listener)
	if err != nil {
		fmt.Fprintf(os.Stderr, "serve: %v", err)
		os.Exit(1)
	}
}

// answer to proposer RPC

func (re *Recorder) HandleESP(req *ProposerMessage) *RecorderResponse {

	var response RecorderResponse

	// send the last decided index details to the proxy, if available
	if len(req.DecidedSlots) > 0 {
		d := Decision{
			indexes:   make([]int, 0),
			decisions: make([][]string, 0),
		}

		for i := 0; i < len(req.DecidedSlots); i++ {
			d.indexes = append(d.indexes, int(req.DecidedSlots[i].Index))
			d.decisions = append(d.decisions, req.DecidedSlots[i].Ids)
		}
		re.recorderToProxyChan <- d
	}

	if len(req.P.ClientBatches) == 0 {
		// if there are only hashes, then check if all the client batches are available in the shared pool. if yes send a positive response depending on the consensus rules
		allBatchesFound := re.findAllBatches(req.P.Ids)
		if !allBatchesFound {
			response.HasClientBacthes = false

		} else {
			// process using the recorder logic
			response.S, response.F, response.M = re.espImpl(req.Index, req.P)
		}
	}

	// if not send a negative response

	// Mark the time of the proposal message for the proposer

	return &response
}

// answer to fetch request

func (r *Recorder) HandleFtech(req *DecideRequest) *DecideResponse {
	var response DecideResponse
	//todo implement
	return &response
}

// check of all the batches are available in the store

func (re *Recorder) findAllBatches(ids []string) bool {
	for i := 0; i < len(ids); i++ {
		_, ok := re.clientBatches.Get(ids[i])
		if !ok {
			return false
		}
	}
	return true
}

// main recorder logic goes here

func (re *Recorder) espImpl(index int64, p *ProposerMessage_Proposal) (int64, *RecorderResponse_Proposal, *RecorderResponse_Proposal) {
	
	re.instanceCreationMutex.Lock()
	
	for int64(len(re.slots)) < index+1 {
		re.slots = append(re.slots, RecorderSlot{
			Mutex: &sync.Mutex{},
			S:     0,
			F:     Value{},
			A:     Value{},
			M:     Value{},
		})
	}
	re.instanceCreationMutex.Unlock()
	
	// from here Oct 26 23.14
	
	//todo
	
	return -1, nil, nil
	
}
