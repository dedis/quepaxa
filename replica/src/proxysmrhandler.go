package raxos

import (
	"fmt"
	"raxos/common"
	"raxos/proto/client"
	"time"
)

// checks if the two string arrays are the same

func (pr *Proxy) hasSameBatches(array1 []string, array2 []string) bool {
	if len(array1) != len(array2) {
		return false
	}
	for i := 0; i < len(array1); i++ {
		if array1[i] != array2[i] {
			return false
		}
	}
	return true
}

// returns true if the decision is same as the proposed value or if I have not proposed anything before

func (pr *Proxy) decidedTheProposedValue(index int, decisions []string) bool {
	if pr.replicatedLog[index].proposedBatch == nil {
		// i have not proposed anything
		return true
	}
	if pr.hasSameBatches(pr.replicatedLog[index].proposedBatch, decisions) {
		return true
	}
	return false
}

// for each item in the list, if it is found in the toBeProposed, then delete it

func (pr *Proxy) removeDecidedItemsFromFutureProposals(items []string) {
	for i := 0; i < len(items); i++ {
		position := -1
		for j := 0; j < len(pr.toBeProposed); j++ {
			if items[i] == pr.toBeProposed[j] {
				position = j
				break
			}
		}
		if position != -1 {
			pr.toBeProposed[position] = pr.toBeProposed[len(pr.toBeProposed)-1]
			pr.toBeProposed = pr.toBeProposed[:len(pr.toBeProposed)-1]
		}
	}
}

// apply the SMR logic for each client request

func (pr *Proxy) applySMRLogic(batch client.ClientBatch) client.ClientBatch {
	//todo implement
	return batch // todo change this later
}

// execute a single client batch

func (pr *Proxy) executeClientBatch(s string) (*client.ClientBatch, bool) {
	batch, ok := pr.clientBatchStore.Get(s)
	if !ok {
		return nil, false
	}
	outputBatch := pr.applySMRLogic(batch)
	return &outputBatch, true
}

// send the client response to client

func (pr *Proxy) sendClientResponse(batches []*client.ClientBatch) {

	for i := 0; i < len(batches); i++ {
		if batches[i].Sender == -1 {
			continue
		}
		pr.sendMessage(batches[i].Sender, common.RPCPair{
			Code: pr.clientBatchRpc,
			Obj:  batches[i],
		})

		pr.debug("proxy sent a client response  "+fmt.Sprintf("%v", batches[i]), -1)
	}
}

// update the state machine by executing all the commands from the committedIndex to len(log)-1
// record the last committed time

func (pr *Proxy) updateStateMachine(sendResponse bool) {
	for i := pr.committedIndex + 1; i < int64(len(pr.replicatedLog)); i++ {
		if pr.replicatedLog[i].decided == true {
			if len(pr.replicatedLog[i].decidedBatch) == 0 {
				panic("should not happen")
			}
			pr.debug("proxy calling update state machine and found a new decided slot  "+fmt.Sprintf("%v", i), 1)

			for j := 0; j < len(pr.replicatedLog[i].decidedBatch); j++ {
				// check if each batch exists
				_, ok := pr.clientBatchStore.Get(pr.replicatedLog[i].decidedBatch[j])
				if !ok {
					pr.proxyToProposerFetchChan <- FetchRequest{ids: pr.replicatedLog[i].decidedBatch}
					pr.debug("proxy cannot commit because the client batches are missing for decided slot  "+fmt.Sprintf("%v", i)+" hence requesting  "+fmt.Sprintf("%v", pr.replicatedLog[i].decidedBatch), 1)
					return
				}
			}

			pr.debug("proxy has all client batches to commit  "+fmt.Sprintf("%v", i), 0)

			var responseBatches []*client.ClientBatch
			for j := 0; j < len(pr.replicatedLog[i].decidedBatch); j++ {
				var responseBatch *client.ClientBatch
				responseBatch, ok := pr.executeClientBatch(pr.replicatedLog[i].decidedBatch[j])
				if !ok {
					panic("did not find the client batch")
				}
				responseBatches = append(responseBatches, responseBatch)
			}
			pr.lastTimeCommitted = time.Now()
			pr.debug("proxy committed  "+fmt.Sprintf("%v", pr.committedIndex+1), 8)
			pr.replicatedLog[i].committed = true
			// empty the proposed batch
			pr.replicatedLog[i].proposedBatch = make([]string, 0)
			pr.committedIndex++
			if sendResponse {
				pr.sendClientResponse(responseBatches)
			}
		} else {
			break
		}
	}
}

// revoke a single instance by proposing the same command proposed before

func (pr *Proxy) revokeInstance(instance int64) {

	pr.debug("proxy revoking instance  "+fmt.Sprintf("%v", pr.replicatedLog[instance]), 9)

	if pr.replicatedLog[instance].decided == true {
		panic("revoking an already decided entry")
	}

	strProposals := pr.replicatedLog[instance].proposedBatch

	if strProposals != nil || len(strProposals) > 0 {
		// I have not proposed for this index before
		panic("should this happen?")
	}
	strProposals = []string{"nil"}

	btchProposals := make([]client.ClientBatch, 0)

	btchProposals = append(btchProposals, client.ClientBatch{
		Sender:   -1,
		Messages: nil,
		Id:       "nil",
	})

	if len(strProposals) != len(btchProposals) {
		panic("lengths do not match")
	}

	newProposalRequest := ProposeRequest{
		instance:             instance,
		proposalStr:          strProposals,
		proposalBtch:         btchProposals,
		msWait:               int(pr.getLeaderWait(pr.getLeaderSequence(instance))),
		lastDecidedIndexes:   nil,
		lastDecidedDecisions: nil,
		leaderSequence:       pr.getLeaderSequence(instance),
	}

	pr.proxyToProposerChan <- newProposalRequest

	pr.debug("proxy revoked instance with new Proposal Request  "+fmt.Sprintf("%v", newProposalRequest), 1)

	pr.replicatedLog[instance] = Slot{
		proposedBatch: strProposals,
		decidedBatch:  pr.replicatedLog[instance].decidedBatch,
		decided:       pr.replicatedLog[instance].decided,
		committed:     pr.replicatedLog[instance].committed,
	}
}

// revoke all the instances from the last committed index to len log

func (pr *Proxy) revokeInstances() {
	for i := pr.committedIndex + 1; i < int64(len(pr.replicatedLog)); i++ {
		if pr.replicatedLog[i].decided == false {
			pr.revokeInstance(i)
		}
	}
}

// handler for propose response from the proposer

func (pr *Proxy) handleProposeResponse(message ProposeResponse) {

	if message.index != -1 && message.decisions != nil {

		pr.debug("proxy received a proposal response from the proxy  "+fmt.Sprintf("%v", message), -1)

		if pr.replicatedLog[message.index].decided == false {
			pr.replicatedLog[message.index].decided = true
			pr.replicatedLog[message.index].decidedBatch = message.decisions

			pr.debug("proxy decided as a result of propose "+fmt.Sprintf(" for instance %v with initial value", message.index, message.decisions[0]), 1)

			if !pr.decidedTheProposedValue(message.index, message.decisions) {
				pr.debug("proxy decided  a different proposal, hence putting back stuff to propose later", 0)
				pr.toBeProposed = append(pr.toBeProposed, pr.replicatedLog[message.index].proposedBatch...)
			}
			// remove the decided batches from toBeProposed
			pr.removeDecidedItemsFromFutureProposals(pr.replicatedLog[message.index].decidedBatch)
		}

		// update SMR -- if all entries are available
		pr.updateStateMachine(true)

		// look at the last time committed, and revoke if needed using no-ops
		if time.Now().Sub(pr.lastTimeCommitted).Milliseconds() > int64(pr.leaderTimeout*100000) { //todo change the revoke timeout
			// revoke all the instances from last committed index
			pr.debug("proxy revoking because it has not committed anything recently  ", 10)
			pr.revokeInstances()
		}

		// add the decided value to proxy's lastDecidedIndexes, lastDecidedDecisions
		pr.lastDecidedIndexes = append(pr.lastDecidedIndexes, message.index)
		pr.lastDecidedDecisions = append(pr.lastDecidedDecisions, message.decisions)

	}
}

// return the highest from the array

func (pr *Proxy) getHighestIndex(indexes []int) int {
	highest := indexes[0]
	for i := 0; i < len(indexes); i++ {
		if indexes[i] > highest {
			highest = indexes[i]
		}
	}
	return highest
}

// mark the entries in the replicated log, and if possible execute

func (pr *Proxy) handleRecorderResponse(message Decision) {

	pr.debug("proxy received decisions from the recorder  "+fmt.Sprintf("%v", message), -1)
	if len(message.indexes) != len(message.decisions) {
		panic("number of decided items and number of decisions do not match")
	}

	highestIndex := pr.getHighestIndex(message.indexes)

	for len(pr.replicatedLog) < int(highestIndex)+1 {
		pr.replicatedLog = append(pr.replicatedLog, Slot{
			proposedBatch: nil,
			decidedBatch:  nil,
			decided:       false,
			committed:     false,
		})
	}

	for i := 0; i < len(message.indexes); i++ {
		index := message.indexes[i]
		batches := message.decisions[i]
		if len(batches) == 0 {
			panic("should not happen")
		}
		if pr.replicatedLog[index].decided == false {
			pr.replicatedLog[index].decided = true
			pr.replicatedLog[index].decidedBatch = batches
			pr.debug("proxy decided from the recorder response "+fmt.Sprintf(" instance %v with batches %v", index, batches), 1)
			if !pr.decidedTheProposedValue(index, batches) {
				pr.toBeProposed = append(pr.toBeProposed, pr.replicatedLog[index].proposedBatch...)
			}
			pr.removeDecidedItemsFromFutureProposals(batches)
		}
	}

	// update SMR -- if all entries are available
	pr.updateStateMachine(false)

}

// save the batch in the store

func (pr *Proxy) handleFetchResponse(response FetchResposne) {
	pr.debug("proxy received fetch response from the proposer "+fmt.Sprintf("%v", response), 1)
	for i := 0; i < len(response.batches); i++ {
		pr.clientBatchStore.Add(response.batches[i])
	}

	pr.debug("proxy update smr after fetch response, note that last committed index is "+fmt.Sprintf("%v", pr.committedIndex), 1)
	if int64(len(pr.replicatedLog)) > pr.committedIndex+1 {
		pr.debug("the state of the next instance is "+fmt.Sprintf("%v", pr.replicatedLog[pr.committedIndex+1]), 1)
	}
	pr.updateStateMachine(true)
}

// return the immutable leader sequence for instance

func (pr *Proxy) getLeaderSequence(instance int64) []int64 {
	if pr.leaderMode == 0 {
		// fixed order
		// assumes that node names start with 1
		rA := make([]int64, 0)

		for i := 0; i < pr.numReplicas; i++ {
			rA = append(rA, int64(i+1))
		}

		return rA
	}
	
	panic("should not happen")
}

// return the pre-agreed, non changing waiting time for the instance by the proposer todo

func (pr *Proxy) getLeaderWait(sequence []int64) int64 {

	for j := 0; j < len(sequence); j++ {
		if sequence[j] == pr.name {
			return pr.leaderTimeout * int64(j)
		}
	}
	panic("should not happen")
}
