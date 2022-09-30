package raxos

// handler for propose response

func (pr *Proxy) handleProposeResponse(message ProposeResponse) {
	if message.uniqueId == pr.replicatedLog[message.index].uniqueId {
		// decided the value I proposed
		if !pr.exec {
			// send back the response to the client if self is who decided that
			//todo
		} else {
			// update SMR -- if all entries are available, if not look at the last time committed, and revoke if needed using no-ops
			// send back the response to client if self is who decided that
			//todo
		}
	} else {
		// re-propose to a new index
		//todo
	}

}

// print the consensus log to file

func (pr *Proxy) printConsensusLog() {
	//todo
}

// mark the entries in the replicated log, and if possible execute

func (pr *Proxy) handleRecorderResponse(message Decision) {
	// todo 
}
