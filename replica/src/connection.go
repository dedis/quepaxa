package raxos

import (
	"context"
	"fmt"
)

type GRPCConnection struct {
	Recorder *Recorder
}

// answer to proposer RPC

func (gc *GRPCConnection) ESP(ctx context.Context, req *ProposerMessage) (*RecorderResponse, error) {
	gc.Recorder.debug("received esp request "+fmt.Sprintf("%v", req), -1)
	var response *RecorderResponse
	response = gc.Recorder.HandleESP(req)
	gc.Recorder.debug("recorder responded to esp request "+fmt.Sprintf("%v", response), -1)
	if response == nil {
		panic("should this happen?")
	}
	return response, nil
}

// for gRPC forward compatibility

func (gc *GRPCConnection) mustEmbedUnimplementedConsensusServer() {
	// no need to implement
}

// answer to fetch Request

func (gc *GRPCConnection) FetchBatches(ctx context.Context, req *DecideRequest) (*DecideResponse, error) {
	gc.Recorder.debug("received fetch batch request "+fmt.Sprintf("%v", req), 0)
	var response *DecideResponse
	response = gc.Recorder.HandleFetch(req)
	gc.Recorder.debug("responded to fetch batch request "+fmt.Sprintf("%v", response), 0)
	return response, nil
}
