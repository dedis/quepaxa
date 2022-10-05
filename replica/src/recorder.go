package raxos

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
	"os"
	"raxos/proto/consensus"
)

type Recorder struct {
	address       string       // address to listen for gRPC connections
	listener      net.Listener // socket for gRPC connections
	server        *grpc.Server // gRPC server
	connection    *consensus.GRPCConnection
	clientBatches *ClientBatchStore
	// pointer to the time array
}

// instantiate a new Recorder

func NewRecorder(address string) *Recorder {

	re := Recorder{
		address: address,
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

	// Mark the time of the propose message for the proposer

	return &response
}
