package raxos

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
	"os"
	"raxos/configuration"
	"raxos/proto/consensus"
	"strconv"
	"sync"
	"time"
)

type Recorder struct {
	address               string       // address to listen for gRPC connections
	listener              net.Listener // socket for gRPC connections
	server                *grpc.Server // gRPC server
	connection            *consensus.GRPCConnection
	clientBatches         *ClientBatchStore
	lastSeenTimeProposers []*time.Time // last seen times of each proposer
	recorderToProxyChan   chan Decision
	name                  int64
	slots                 []RecorderSlot // recorder side replicated log
	cfg                   configuration.InstanceConfig
}

type RecorderSlot struct {
	mutex *sync.Mutex
	//todo
}

// instantiate a new Recorder

func NewRecorder(cfg configuration.InstanceConfig, clientBatches *ClientBatchStore, lastSeenTimeProposers []*time.Time, recorderToProxyChan chan Decision, name int64) *Recorder {

	re := Recorder{
		address:               "",
		clientBatches:         clientBatches,
		lastSeenTimeProposers: lastSeenTimeProposers,
		recorderToProxyChan:   recorderToProxyChan,
		name:                  name,
		slots:                 make([]RecorderSlot, 0),
		cfg:                   cfg,
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
	r.connection = &consensus.GRPCConnection{Recorder: r}
	consensus.RegisterConsensusServer(r.server, r.connection)

	// start listener
	listener, err := net.Listen("tcp", r.address)
	r.listener = listener
	if err != nil {
		fmt.Fprintf(os.Stderr, "listen: %v", err)
		os.Exit(1)
	}
	err = r.server.Serve(listener)
	if err != nil {
		fmt.Fprintf(os.Stderr, "serve: %v", err)
		os.Exit(1)
	}
}

// answer to proposer RPC

func (re *Recorder) HandleESP(req *consensus.ProposerMessage) *consensus.RecorderResponse {

	var response consensus.RecorderResponse

	// todo

	// send the last decided index details to the proxy, if available

	// if there are only hashes, then check if all the client batches are available in the shared pool. if yes send a positive response depending on the consensus rules

	// if not send a negative response

	// Mark the time of the proposal message for the proposer

	return &response
}

// answer to fetch request

func (r *Recorder) HandleFtech(req *consensus.DecideRequest) *consensus.DecideResponse {
	var response consensus.DecideResponse
	//todo implement
	return &response
}
