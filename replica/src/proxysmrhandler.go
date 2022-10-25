package raxos

import (
	"raxos/common"
	"raxos/proto/client"
	"time"
)

// handler for propose response

func (pr *Proxy) handleProposeResponse(message ProposeResponse) {
	if pr.decidedTheProposedValue(message) {
		if pr.replicatedLog[message.index].decided == false {
			// decide the value I proposed
			pr.replicatedLog[message.index].decided = true
			pr.replicatedLog[message.index].decidedBatch = message.decisions

			// update SMR -- if all entries are available
			pr.updateStateMachine()
			
			// from here: if not look at the last time committed, and revoke if needed using no-ops

			// send back the response to client if self is who decided that
			//todo
		}
	} else {
		// re-propose to a new index
		//todo
	}
	// add the decided value to proxy's lastDecidedIndexes, lastDecidedDecisions and lastDecidedUniqueIds

}

// update the state machine by executing all the commands from the committedIndex to len(log)-1
// record the last committed time

func (pr *Proxy) updateStateMachine() {
	for i := pr.committedIndex + 1; i < int64(len(pr.replicatedLog)); i++ {

		if pr.replicatedLog[i].decided == true {

			for j := 0; j < len(pr.replicatedLog[i].decidedBatch); j++ {
				var responseBatch *client.ClientBatch
				responseBatch = pr.executeClientBatch(pr.replicatedLog[i].decidedBatch[j])
				pr.sendClientRespose(responseBatch)
				pr.committedIndex++
				pr.lastTimeCommitted = time.Now()
			}
		} else {
			break
		}
	}
}

// execute a single client batch

func (pr *Proxy) executeClientBatch(s string) *client.ClientBatch {
	batch, ok := pr.clientBatchStore.Get(s)
	if !ok {
		panic("decided batch not in the store")
	}
	// todo put application logic here
	var outputBatch client.ClientBatch
	outputBatch = pr.applySMRLogic(batch)
	return &outputBatch
}

// apply the SMR logic for each client request

func (pr *Proxy) applySMRLogic(batch client.ClientBatch) client.ClientBatch {
	//todo implement
	return batch // todo change this later
}

// send the client response to client

func (pr *Proxy) sendClientRespose(batch *client.ClientBatch) {
	pr.sendMessage(batch.Sender, common.RPCPair{
		Code: pr.clientBatchRpc,
		Obj:  batch,
	})
}

// mark the entries in the replicated log, and if possible execute

func (pr *Proxy) handleRecorderResponse(message Decision) {
	// todo
}

// returns true if the decision is same as the proposed value

func (pr *Proxy) decidedTheProposedValue(message ProposeResponse) bool {
	index := message.index
	if pr.hasSameBatches(pr.replicatedLog[index].proposedBatch, message.decisions) {
		return true
	}
	return false
}

// checks if the two arrays are same

func (pr *Proxy) hasSameBatches(batch []string, decisions []string) bool {
	if len(batch) != len(decisions) {
		return false
	}
	for i := 0; i < len(batch); i++ {
		if batch[i] != decisions[i] {
			return false
		}
	}
	return true
}
