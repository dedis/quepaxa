package raxos

import (
	"context"
	"raxos/proto/consensus"
)

type Recorder struct {
	// gRPC listener variables
	// pointer to the time array
}

// instantiate a new Recorder

func NewRecorder() *Recorder {

	re := Recorder{}

	return &re
}

// start listening to gRPC connection

func (r *Recorder) NetworkInit() {

}

// answer to proposer RPC

func (re *Recorder) ESP(ctx context.Context, req *consensus.ProposerMessage) (*consensus.RecorderResponse, error) {

	var response consensus.RecorderResponse
	
	// todo
	
	// send the last decided index details to the proxy
	
	// if there are only hashes, then check if all the client batches are available in the shared pool. if yes send a positive response depending on the consensus rules
	
	// if not send a negative response
	
	// Mark the time of the propose message for the proposer
	
	return &response, nil
}
