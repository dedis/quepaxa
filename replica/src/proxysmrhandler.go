package raxos

import (
	"raxos/common"
	"raxos/proto/client"
	"strconv"
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

// remove the items in array from pr.toBeProposed

func (pr *Proxy) removeDecidedItemsFromFutureProposals(array []string) {

	// Create a set to store the elements of array
	set := make(map[string]bool)
	for _, elem := range array {
		set[elem] = true
	}

	// Remove elements from pr.toBeProposed if they exist in the set
	j := 0
	for i := 0; i < len(pr.toBeProposed); i++ {
		if !set[pr.toBeProposed[i]] {
			pr.toBeProposed[j] = pr.toBeProposed[i]
			j++
		}
	}

	// Truncate pr.toBeProposed to remove the remaining elements
	pr.toBeProposed = pr.toBeProposed[:j]
}

// apply the SMR logic for client requests

func (pr *Proxy) applySMRLogic(batches []client.ClientBatch) []client.ClientBatch {
	responses := pr.benchmark.Execute(batches)
	return responses
}

// execute a client batches

func (pr *Proxy) executeClientBatches(s []string) []client.ClientBatch {
	batches := make([]client.ClientBatch, 0)
	for i := 0; i < len(s); i++ {
		batch, ok := pr.clientBatchStore.Get(s[i])
		if !ok {
			panic("should not happen")
		}
		batches = append(batches, batch)
	}

	outputBatch := pr.applySMRLogic(batches)
	return outputBatch
}

// send the client response to client

func (pr *Proxy) sendClientResponse(batches []client.ClientBatch) {

	for i := 0; i < len(batches); i++ {
		if batches[i].Sender == -1 {
			continue
		}
		pr.sendMessage(batches[i].Sender, common.RPCPair{
			Code: pr.clientBatchRpc,
			Obj:  &batches[i],
		})

		pr.debug("proxy sent a client response for batch id "+batches[i].Id, 0)
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
			pr.debug("proxy calling update state machine and found a new decided slot  "+strconv.Itoa(int(i)), 0)

			for j := 0; j < len(pr.replicatedLog[i].decidedBatch); j++ {
				// check if each batch exists
				_, ok := pr.clientBatchStore.Get(pr.replicatedLog[i].decidedBatch[j])
				if !ok {
					pr.proxyToProposerFetchChan <- FetchRequest{ids: pr.replicatedLog[i].decidedBatch}
					pr.debug("proxy cannot commit because the client batches are missing for decided slot  "+strconv.Itoa(int(i))+" hence requesting  "+pr.replicatedLog[i].decidedBatch[j], 0)
					return
				}
			}

			pr.debug("proxy has all client batches to commit slot "+strconv.Itoa(int(i)), 0)

			responseBatches := pr.executeClientBatches(pr.replicatedLog[i].decidedBatch)

			pr.lastTimeCommitted = time.Now()
			pr.debug("proxy committed  index "+strconv.Itoa(int(pr.committedIndex+1)), 0)
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

	//look at the last time committed, and revoke if needed using no-ops
	if time.Now().Sub(pr.lastTimeCommitted).Microseconds() > int64(pr.leaderTimeout*2*int64(pr.numReplicas)) {
		//revoke all the instances from last committed index
		pr.debug("proxy revoking because it has not committed anything recently  ", 20)
		//pr.revokeInstances()
	}
}

// revoke a single instance by proposing the same command proposed before

func (pr *Proxy) revokeInstance(instance int64) {

	//pr.debug("proxy revoking instance  "+fmt.Sprintf("%v", pr.replicatedLog[instance]), 9)

	if pr.replicatedLog[instance].decided == true {
		panic("revoking an already decided entry")
	}

	strProposals := pr.replicatedLog[instance].proposedBatch

	if strProposals != nil || len(strProposals) > 0 {
		// I have proposed for this index before
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
		isLeader:             false,
		lastDecidedIndexes:   nil,
		lastDecidedDecisions: nil,
	}

	pr.proxyToProposerChan <- newProposalRequest

	//pr.debug("proxy revoked instance with new Proposal Request  "+fmt.Sprintf("%v", newProposalRequest), 1)

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

		//pr.debug("proxy received a proposal response from the proxy  "+fmt.Sprintf("%v", message), -1)

		if pr.replicatedLog[message.index].decided == false {
			pr.replicatedLog[message.index].decided = true
			pr.replicatedLog[message.index].decidedBatch = message.decisions
			pr.replicatedLog[message.index].proposer = message.proposer
			pr.replicatedLog[message.index].s = message.s

			//pr.debug("proxy decided as a result of propose "+fmt.Sprintf(" for instance %v with initial value", message.index, message.decisions[0]), 20)
			pr.updateEpochTime(message.index)

			if !pr.decidedTheProposedValue(message.index, message.decisions) {
				//pr.debug("proxy decided  a different proposal, hence putting back stuff to propose later", 0)
				pr.toBeProposed = append(pr.toBeProposed, pr.replicatedLog[message.index].proposedBatch...)
				pr.replicatedLog[message.index].proposedBatch = nil
			}

			pr.removeDecidedItemsFromFutureProposals(pr.replicatedLog[message.index].decidedBatch)

		}

		// update SMR -- if all entries are available
		pr.updateStateMachine(true)

		// add the decided value to proxy's lastDecidedIndexes, lastDecidedDecisions
		pr.lastDecidedIndexes = append(pr.lastDecidedIndexes, message.index)
		pr.lastDecidedProposers = append(pr.lastDecidedProposers, message.proposer)
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

	//pr.debug("proxy received decisions from the recorder  "+fmt.Sprintf("%v", message), 11)
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
			s:             0,
		})
	}

	for i := 0; i < len(message.indexes); i++ {
		index := message.indexes[i]
		batches := message.decisions[i]
		proposer := message.proposers[i]
		if len(batches) == 0 {
			panic("should not happen")
		}
		if pr.replicatedLog[index].decided == false {
			pr.replicatedLog[index].decided = true
			pr.replicatedLog[index].decidedBatch = batches
			pr.replicatedLog[index].proposer = proposer

			pr.updateEpochTime(index)
			pr.debug("proxy decided from the recorder response "+strconv.Itoa(index), 20)
			if !pr.decidedTheProposedValue(index, batches) {
				pr.toBeProposed = append(pr.toBeProposed, pr.replicatedLog[index].proposedBatch...)
				pr.replicatedLog[index].proposedBatch = nil
			}
			pr.removeDecidedItemsFromFutureProposals(batches)
		}
	}

	// update SMR -- if all entries are available
	pr.updateStateMachine(false)

}

// save the batch in the store

func (pr *Proxy) handleFetchResponse(response FetchResposne) {
	//pr.debug("proxy received fetch response from the proposer "+fmt.Sprintf("%v", response), 1)
	for i := 0; i < len(response.batches); i++ {
		pr.clientBatchStore.Add(response.batches[i])
	}
	//pr.debug("proxy update smr after fetch response, note that last committed index is "+fmt.Sprintf("%v", pr.committedIndex), 1)
	pr.updateStateMachine(true)
}

// send decisions to the proposer

func (pr *Proxy) handleDecisionNotification() {

	if len(pr.lastDecidedIndexes) == 0 {
		return
	}
	newDecision := Decision{
		indexes:   pr.lastDecidedIndexes,
		decisions: pr.lastDecidedDecisions,
		proposers: pr.lastDecidedProposers,
	}

	pr.proxyToProposerDecisionChan <- newDecision
	//pr.debug("proxy sent a decisions to proposer  "+fmt.Sprintf("%v", newDecision), 11)

	// reset the variables
	pr.lastDecidedIndexes = make([]int, 0)
	pr.lastDecidedProposers = make([]int32, 0)
	pr.lastDecidedDecisions = make([][]string, 0)
}
