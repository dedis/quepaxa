package raxos

import (
	"raxos/configuration"
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
	//todo
}

// ProposeResponse is the message type sent from proposer to proxy

type ProposeResponse struct {
	//todo
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

/*
	create a new server instance, inside which there are proxy instance, proposer instances and recorder instance. initialize all fields
*/

func New(cfg *configuration.InstanceConfig, name int64, logFilePath string, batchSize int64, batchTime int64, leaderTimeout int64, pipelineLength int64, benchmark int64, debugOn bool, debugLevel int) *Server {

	sr := Server{}

	return &sr
}
