package raxos

type Proposer struct {
	peers []peer // gRPC connection list
	// chan from proxy
	// chan to proxy
	// pointer to the time array
}

// instantiate a new Proposer

func NewProposer(peers []peer) *Proposer {

	pr := Proposer{
		peers: peers,
	}

	return &pr
}

// infinte loop listening to the server channel

func (prop *Proposer) runProposer() {
	go func() {
		for true {
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
