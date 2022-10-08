package raxos

import "time"

type Proposer struct {
	name int64
	threadId int64
	peers []peer // gRPC connection list
	proxyToProposerChan chan ProposeRequest
	proposerToProxyChan chan ProposeResponse
	lastSeenTimes []*time.Time 
}

// instantiate a new Proposer

func NewProposer(name int64, threadId int64, peers []peer, proxyToProposerChan chan ProposeRequest, proposerToProxyChan chan ProposeResponse, lastSeenTimes []*time.Time ) *Proposer {

	pr := Proposer{
		name:                name,
		threadId:            threadId,
		peers:               peers,
		proxyToProposerChan: proxyToProposerChan,
		proposerToProxyChan: proposerToProxyChan,
		lastSeenTimes:       lastSeenTimes,
	}

	return &pr
}

// infinite loop listening to the server channel

func (prop *Proposer) runProposer() {
	go func() {
		for true {
			//todo
			// get a new request
			// wait for the time to pass
			// propose for the fast path
			// wait for the responses
			// if needed start the slow path
			// run to completion
			// send the response back to the proxy
		}
	}()
}
