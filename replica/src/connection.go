package raxos

import (
	"context"
)

// this file defines the gRPC receive side code for the recorder. the actual Recorder implementation is in the recorder.go

// GRPCConnection is a grpc wrapper for recorder

type GRPCConnection struct {
	Recorder *Recorder
}

func (gc *GRPCConnection) InformDecision(ctx context.Context, decisions *Decisions) (*Empty, error) {
	//gc.Recorder.debug("received decisions response", 0)
	var response *Empty
	response = &Empty{}
	gc.Recorder.HandleDecisions(decisions)
	return response, nil
}

// answer to proposer RPC

func (gc *GRPCConnection) ESP(ctx context.Context, req *ProposerMessage) (*RecorderResponse, error) {
	//gc.Recorder.debug("received esp request from"+strconv.Itoa(int(req.Sender)), 0)
	var response *RecorderResponse
	response = gc.Recorder.HandleESP(req)
	//gc.Recorder.debug("recorder responded to esp request from "+strconv.Itoa(int(req.Sender)), 0)
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
	//gc.Recorder.debug("received fetch batch request", 0)
	var response *DecideResponse
	response = gc.Recorder.HandleFetch(req)
	//gc.Recorder.debug("responded to fetch batch request ", 0)
	return response, nil
}
