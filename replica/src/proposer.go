package raxos

type Proposer struct {
	// gRPC connection list
	// chan from server
	// chan to server
	// pointer to the time array
}

// instantiate a new Proposer

func NewProposer() *Proposer {

	pr := Proposer{}

	return &pr
}

// infinte loop listening to the server channel

func ( prop *Proposer) runProposer() {
	go func(){
		for true{
			// get a new request
			// propose for the fast path
			// wait for the responses
			// if needed start the slow path
			// run to completion 
			// send the response back to the proxy
		} 
	}()
}
