package consensus

import (
	"context"
	raxos "raxos/replica/src"
)

type GRPCConnection struct {
	Recorder *raxos.Recorder
}

// answer to proposer RPC

func (gc *GRPCConnection) ESP(ctx context.Context, req *ProposerMessage) (*RecorderResponse, error) {

	var response *RecorderResponse
	response = gc.Recorder.HandleESP(req)
	return response, nil
}

// for gRPC forward compatibility

func (gc *GRPCConnection) mustEmbedUnimplementedConsensusServer() {
	// no need to implement
}

// answer to fetch Request

func (gc *GRPCConnection) Fetch(ctx context.Context, req *DecideRequest) (*DecideResponse, error) {

	var response *DecideResponse
	response = gc.Recorder.HandleFtech(req)
	return response, nil
}