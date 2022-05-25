package raxos

import (
	"fmt"
	"math"
	"raxos/proto"
	"strconv"
	"strings"
)

/*
	mark the slot as decided and inform the upper layer
*/

func (in *Instance) recordProposerDecide(consensusMessage *proto.GenericConsensus) {
	in.proposerReplicatedLog[consensusMessage.Index].S = consensusMessage.S
	in.proposerReplicatedLog[consensusMessage.Index].decided = consensusMessage.D
	in.proposerReplicatedLog[consensusMessage.Index].decision = consensusMessage.DS
	in.proposerReplicatedLog[consensusMessage.Index].proposer = consensusMessage.PR

	in.proposerReplicatedLog[consensusMessage.Index].P = []*proto.GenericConsensusValue{}
	in.proposerReplicatedLog[consensusMessage.Index].E = []*proto.GenericConsensusValue{}
	in.proposerReplicatedLog[consensusMessage.Index].C = []*proto.GenericConsensusValue{}
	in.proposerReplicatedLog[consensusMessage.Index].proposeResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].spreadEResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].spreadCGatherEResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].gatherCResponses = []*proto.GenericConsensus{}

	if consensusMessage.Index == in.lastDecidedIndex+1 {
		in.awaitingDecision = false
	}
	in.updateLastDecidedIndex()
}

/*
	update the last decided index and call the SMR
*/

func (in *Instance) updateLastDecidedIndex() {
	for i := in.lastDecidedIndex + 1; i < int64(len(in.proposerReplicatedLog)); i++ {
		if in.proposerReplicatedLog[i].decided == true {
			in.lastDecidedIndex++
			in.debug("Updated last decided index "+strconv.Itoa(int(in.lastDecidedIndex)), 1)
		} else {
			break
		}
	}
	in.updateStateMachine()
}

/*
	invoked when last decided index is updated
	marks the slot as committed
	sends the response back to the client
*/

func (in *Instance) updateStateMachine() { // from here
	for i := in.committedIndex + 1; i < in.lastDecidedIndex+1; i++ {

		decision := in.proposerReplicatedLog[i].decision.Id
		decision = strings.Split(decision, "-")[0]
		// if the decision is a sequence of hashes, then invoke the following for each hash: check if all the hashes exist and succeed only if everything exists
		blocks := make([]*proto.MessageBlock, 0)
		hashes := strings.Split(decision, ":")

		// remove white spaces of each hash
		for j := 0; j < len(hashes); j++ {
			hashes[j] = strings.TrimSpace(hashes[j])
			in.debug("hash "+strconv.Itoa(j)+" is "+hashes[j], 0)
		}

		// for each hash, check if the message block exists
		for j := 0; j < len(hashes); j++ {
			if len(hashes[j]) == 0 {
				continue
			}
			messageBlock, ok := in.messageStore.Get(hashes[j])
			if !ok {
				in.debug("Couldn't find hash "+hashes[j], 0)
				/*
					I haven't received the actual message block still, request it
				*/
				in.sendMessageBlockRequest(hashes[j])
				return
			}
			blocks = append(blocks, messageBlock)
		}

		// execute each mem block
		for j := 0; j < len(blocks); j++ {
			in.executeAndSendResponse(blocks[j])
		}

		in.proposerReplicatedLog[i].committed = true
		in.committedIndex++
		in.debug("committing instance "+strconv.Itoa(int(i))+" with decision "+decision, 4)
	}
}

/*
	execute each command in the block, get the response array and send the responses back to the client if I am the proposer who created this block
*/

func (in *Instance) executeAndSendResponse(block *proto.MessageBlock) {
	for i := 0; i < len(block.Requests); i++ {
		clientRequestBatch := block.Requests[i]
		var replies []*proto.ClientResponseBatch_SingleClientResponse
		for j := 0; j < len(clientRequestBatch.Requests); j++ {
			//todo add smr invocation here
			replies = append(replies, &proto.ClientResponseBatch_SingleClientResponse{Message: clientRequestBatch.Requests[j].Message})
		}

		// for each client batch, send a response batch if I originated this block

		if block.Sender == in.nodeName {

			responseBatch := proto.ClientResponseBatch{
				Sender:    in.nodeName,
				Receiver:  clientRequestBatch.Sender,
				Responses: replies,
				Id:        clientRequestBatch.Id,
			}
			rpcPair := RPCPair{
				Code: in.clientResponseBatchRpc,
				Obj:  &responseBatch,
			}

			in.sendMessage(clientRequestBatch.Sender, rpcPair)
			in.debug("Sent client response batch to "+strconv.Itoa(int(clientRequestBatch.Sender)), 1)
		}

	}
}

/*
	upon receiving a consensus request, the leader node will propose it to the consensus layer if the proposer is not waiting to decide the current instance
*/

func (in *Instance) handleConsensusRequest(request *proto.ConsensusRequest) {
	if len([]rune(request.Hash)) > 0 {
		if in.nodeName == in.getDeterministicLeader1() { //todo change this for optimizations
			if !in.awaitingDecision {
				in.propose(in.lastDecidedIndex+1, request.Hash)
				in.awaitingDecision = true
				in.debug("Proposed "+request.Hash+" to "+strconv.Itoa(int(in.lastDecidedIndex+1)), 0)
			} else {
				for strings.HasPrefix(request.Hash, ":") {
					request.Hash = request.Hash[1:]
				}

				in.hashProposalsIn <- request.Hash
			}
		}
	}
}

/*
	propose a new block for the slot index
*/

func (in *Instance) propose(index int64, hash string) {
	in.proposerReplicatedLog = in.initializeSlot(in.proposerReplicatedLog, index) // create the slot if not already created
	in.proposerReplicatedLog[index].S = 4
	in.proposerReplicatedLog[index].P = make([]*proto.GenericConsensusValue, 1)
	in.proposerReplicatedLog[index].P[0] = &proto.GenericConsensusValue{
		Id:  hash + "-" + strconv.Itoa(int(in.nodeName)) + "." + strconv.Itoa(in.proposeCounter), // for unique proposer id
		Fit: "",
	}
	in.proposerReplicatedLog[index].proposedValue = hash + "-" + strconv.Itoa(int(in.nodeName)) + "." + strconv.Itoa(in.proposeCounter)
	in.proposeCounter++

	in.proposerReplicatedLog[index].E = make([]*proto.GenericConsensusValue, 0)
	in.proposerReplicatedLog[index].C = make([]*proto.GenericConsensusValue, 0)
	in.proposerReplicatedLog[index].proposeResponses = make([]*proto.GenericConsensus, 0)
	in.proposerReplicatedLog[index].spreadEResponses = make([]*proto.GenericConsensus, 0)
	in.proposerReplicatedLog[index].spreadCGatherEResponses = make([]*proto.GenericConsensus, 0)
	in.proposerReplicatedLog[index].gatherCResponses = make([]*proto.GenericConsensus, 0)

	// send the proposal message to all the recorders
	in.proposerSendMessage(index)
	// reset the P set
	in.proposerReplicatedLog[index].P = make([]*proto.GenericConsensusValue, 0)
}

/*
	broadcast a message to all recorders
*/

func (in *Instance) proposerSendMessage(index int64) {

	destination := in.consensusMessageRecorderDestination

	if in.proposerReplicatedLog[index].decided == true {
		destination = in.consensusMessageCommonDestination
	}
	for i := int64(0); i < in.numReplicas; i++ {

		consensusPropose := proto.GenericConsensus{
			Sender:      in.nodeName,
			Receiver:    i,
			Index:       index,
			S:           in.proposerReplicatedLog[index].S,
			P:           in.removeEmptyValues(in.proposerReplicatedLog[index].P),
			E:           in.removeEmptyValues(in.proposerReplicatedLog[index].E),
			C:           in.removeEmptyValues(in.proposerReplicatedLog[index].C),
			D:           in.proposerReplicatedLog[index].decided,
			DS:          in.proposerReplicatedLog[index].decision,
			PR:          in.proposerReplicatedLog[index].proposer,
			Destination: destination,
		}

		rpcPair := RPCPair{
			Code: in.genericConsensusRpc,
			Obj:  &consensusPropose,
		}
		in.debug("sending a consensus message to "+strconv.Itoa(int(i))+" to destination "+strconv.Itoa(int(destination)), 1)

		in.sendMessage(i, rpcPair)
	}
}

/*
	handler for proposer consensus messages
*/

func (in *Instance) handleProposerConsensusMessage(consensusMessage *proto.GenericConsensus) {
	// since we do not use pipelining,
	if consensusMessage.Index == in.lastDecidedIndex+1 {

		// case 2: a propose reply message from a recorder
		if in.proposerReplicatedLog[consensusMessage.Index].S == consensusMessage.S && consensusMessage.S%4 == 0 {

			in.proposerReplicatedLog[consensusMessage.Index].proposeResponses = append(in.proposerReplicatedLog[consensusMessage.Index].proposeResponses, consensusMessage)
			in.proposerReplicatedLog[consensusMessage.Index].P = in.setUnionProtoValue(in.proposerReplicatedLog[consensusMessage.Index].P, consensusMessage.P[0])
			in.proposerReplicatedLog[consensusMessage.Index].P = in.removeEmptyValues(in.proposerReplicatedLog[consensusMessage.Index].P)
			in.debug("Received propose response from recorder has id = "+consensusMessage.P[0].Id+" and fit = "+consensusMessage.P[0].Fit, 0)

			// case 2.1: on receiving quorum of propose replies
			if int64(len(in.proposerReplicatedLog[consensusMessage.Index].proposeResponses)) == in.numReplicas/2+1 {

				majorityValue := in.getProposalWithMajorityHi(in.proposerReplicatedLog[consensusMessage.Index].proposeResponses)

				// case 2.1.1 a majority of the proposals responses have Hi and I am not decided yet
				if majorityValue.Id != "" && majorityValue.Fit != "" && in.proposerReplicatedLog[consensusMessage.Index].decided == false {
					in.proposerReceivedMajorityProposeWithHi(consensusMessage, majorityValue)
					in.debug("decided in the fast path "+fmt.Sprintf(" %v ", consensusMessage.Index), 4)
					return
				}

				// case 2.1.2 no proposal with majority Hi
				if majorityValue.Id == "" && majorityValue.Fit == "" {
					in.debug("no proposal with majority high "+fmt.Sprintf(" %v ", in.proposerReplicatedLog[consensusMessage.Index]), 4)
					in.proposerReplicatedLog[consensusMessage.Index].S = in.proposerReplicatedLog[consensusMessage.Index].S + 1
					in.proposerSendMessage(consensusMessage.Index)
					return
				}

			}
			return
		}

		// case 3: a spreadE response
		if in.proposerReplicatedLog[consensusMessage.Index].S == consensusMessage.S && in.proposerReplicatedLog[consensusMessage.Index].S%4 == 1 && in.setContainsP(consensusMessage.Index, consensusMessage.E) {
			in.proposerReplicatedLog[consensusMessage.Index].spreadEResponses = append(in.proposerReplicatedLog[consensusMessage.Index].spreadEResponses, consensusMessage)

			// case 3.1 upon receiving a f+1 spreadE message
			if int64(len(in.proposerReplicatedLog[consensusMessage.Index].spreadEResponses)) == in.numReplicas/2+1 {
				in.proposerReplicatedLog[consensusMessage.Index].S = in.proposerReplicatedLog[consensusMessage.Index].S + 1
				in.proposerSendMessage(consensusMessage.Index)
				return

			}
			return

		}

		// case 4: a SpreadCGatherE response
		if in.proposerReplicatedLog[consensusMessage.Index].S == consensusMessage.S && in.proposerReplicatedLog[consensusMessage.Index].S%4 == 2 && in.setContainsP(consensusMessage.Index, consensusMessage.C) {
			in.proposerReplicatedLog[consensusMessage.Index].spreadCGatherEResponses = append(in.proposerReplicatedLog[consensusMessage.Index].spreadCGatherEResponses, consensusMessage)

			// case 4.1: upon receiving a majority SpreadCGatherE responses

			if len(in.proposerReplicatedLog[consensusMessage.Index].spreadCGatherEResponses) == int(in.numReplicas)/2+1 {
				in.proposerReplicatedLog[consensusMessage.Index].S = in.proposerReplicatedLog[consensusMessage.Index].S + 1
				// E ← union of {E’} from all replies in the quorum
				in.proposerReplicatedLog[consensusMessage.Index].E = in.getESetFromESets(in.proposerReplicatedLog[consensusMessage.Index].spreadCGatherEResponses)

				P_Best := in.getBestProposal(in.proposerReplicatedLog[consensusMessage.Index].P)
				E_Best := in.getBestProposal(in.proposerReplicatedLog[consensusMessage.Index].E)

				// case 4.1.1 If best(P) = best(E)
				if P_Best.Id != "" && P_Best.Fit != "" && E_Best.Id != "" && E_Best.Fit != "" {
					if P_Best.Id == E_Best.Id && P_Best.Fit == E_Best.Fit {
						in.debug("decided "+strconv.Itoa(int(consensusMessage.Index))+" in the SpreadCGatherE phase, step is "+strconv.Itoa(int(in.proposerReplicatedLog[consensusMessage.Index].S))+" size of E set is "+strconv.Itoa(len(in.proposerReplicatedLog[consensusMessage.Index].E)), 3)
						in.proposerReceivedMajorityProposeWithHi(consensusMessage, P_Best)
						return
					}
				}

				in.proposerSendMessage(consensusMessage.Index)
				return

			}
			return
		}

		// case 5: a GatherC response
		if in.proposerReplicatedLog[consensusMessage.Index].S == consensusMessage.S && in.proposerReplicatedLog[consensusMessage.Index].S%4 == 3 {
			// update the C set
			in.proposerReplicatedLog[consensusMessage.Index].C = in.setUnionProtoValues(in.proposerReplicatedLog[consensusMessage.Index].C, consensusMessage.C)
			in.proposerReplicatedLog[consensusMessage.Index].C = in.removeEmptyValues(in.proposerReplicatedLog[consensusMessage.Index].C)
			in.proposerReplicatedLog[consensusMessage.Index].gatherCResponses = append(in.proposerReplicatedLog[consensusMessage.Index].gatherCResponses, consensusMessage)

			// case 5.1: upon receiving majority GatherC

			if len(in.proposerReplicatedLog[consensusMessage.Index].gatherCResponses) == int(in.numReplicas)/2+1 {
				in.proposerHandleGatherCQQuorum(consensusMessage)
				in.proposerSendMessage(consensusMessage.Index)
				return
			}
			return

		}

		// case 5: a catch-up message
		if consensusMessage.S > in.proposerReplicatedLog[consensusMessage.Index].S {
			in.consensusCatchUp(consensusMessage)
		}
	}
}

/*
	scan the propose replies set from recorders and assign the number of times the Hi priority appears for each proposal
*/

func (in *Instance) getProposalWithMajorityHi(e []*proto.GenericConsensus) *proto.GenericConsensusValue {
	var ValueCount map[string]int64 // maps each id to hi count
	ValueCount = make(map[string]int64)
	for i := 0; i < len(e); i++ {
		if strings.HasPrefix(e[i].P[0].Fit, strconv.FormatInt(in.Hi, 10)) {
			_, ok := ValueCount[e[i].P[0].Id] // note that this is a unique key
			if !ok {
				ValueCount[e[i].P[0].Id] = 1 // note that each proposal has a unique id
			} else {
				ValueCount[e[i].P[0].Id] = ValueCount[e[i].P[0].Id] + 1
			}
		}
	}

	// check whether there is a value who has f+1 Hi count

	for key, element := range ValueCount {
		if element >= in.numReplicas/2+1 {
			for i := 0; i < len(e); i++ {
				if e[i].P[0].Id == key {
					return e[i].P[0]
				}
			}
		}
	}

	in.debug(fmt.Sprintf("\n%v\n", ValueCount), 4)

	return &proto.GenericConsensusValue{
		Id:  "",
		Fit: "",
	} // there is no proposal with f+1 Hi

}

/*
	called when the proposer receives f+1 propose responses, and each response corresponds to Hi priority for the same value
*/

func (in *Instance) proposerReceivedMajorityProposeWithHi(consensusMessage *proto.GenericConsensus, majorityValue *proto.GenericConsensusValue) {
	in.proposerReplicatedLog[consensusMessage.Index].S = math.MaxInt64
	in.proposerReplicatedLog[consensusMessage.Index].decided = true
	in.proposerReplicatedLog[consensusMessage.Index].decision = majorityValue
	in.proposerReplicatedLog[consensusMessage.Index].proposer = in.getProposer(majorityValue)

	if consensusMessage.Index == in.lastDecidedIndex+1 {
		in.awaitingDecision = false
	}
	in.updateLastDecidedIndex()

	// check what was proposed

	if in.proposerReplicatedLog[consensusMessage.Index].proposedValue != "" && in.proposerReplicatedLog[consensusMessage.Index].proposedValue != majorityValue.Id {
		// the decided item is different from what I proposed: retry proposing
		in.debug("decided a different value than proposed", 4)
		proposedValue := strings.Split(in.proposerReplicatedLog[consensusMessage.Index].proposedValue, "-")[0]
		for strings.HasPrefix(proposedValue, ":") {
			proposedValue = proposedValue[1:]
		}
		in.hashProposalsIn <- proposedValue

	}

	// clean the memory
	in.proposerReplicatedLog[consensusMessage.Index].proposedValue = ""
	in.proposerReplicatedLog[consensusMessage.Index].P = []*proto.GenericConsensusValue{}
	in.proposerReplicatedLog[consensusMessage.Index].E = []*proto.GenericConsensusValue{}
	in.proposerReplicatedLog[consensusMessage.Index].C = []*proto.GenericConsensusValue{}
	in.proposerReplicatedLog[consensusMessage.Index].proposeResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].spreadEResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].spreadCGatherEResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].gatherCResponses = []*proto.GenericConsensus{}

	// send a decide message

	for i := int64(0); i < in.numReplicas; i++ {

		consensusDecide := proto.GenericConsensus{
			Sender:      in.nodeName,
			Receiver:    i,
			Index:       consensusMessage.Index,
			S:           -1, // not needed
			P:           nil,
			E:           nil,
			C:           nil,
			D:           in.proposerReplicatedLog[consensusMessage.Index].decided,
			DS:          in.proposerReplicatedLog[consensusMessage.Index].decision,
			PR:          in.proposerReplicatedLog[consensusMessage.Index].proposer,
			Destination: in.consensusMessageCommonDestination, // to be consumed by both proposer and recorder
		}

		rpcPair := RPCPair{
			Code: in.genericConsensusRpc,
			Obj:  &consensusDecide,
		}
		in.debug("sending a decide consensus message to "+strconv.Itoa(int(i)), 1)

		in.sendMessage(i, rpcPair)
	}

}

/*
	returns the proposer who proposed this value, extracted using the fitness
*/

func (in *Instance) getProposer(value *proto.GenericConsensusValue) int64 {
	priority := value.Fit
	proposer, err := strconv.Atoi(strings.Split(priority, ".")[1])
	if err == nil {
		return int64(proposer)
	} else {
		return -1
	}
}

/*
	check if the set contains the P set of this proposer
*/

func (in *Instance) setContainsP(index int64, e []*proto.GenericConsensusValue) bool {
	pSet := in.proposerReplicatedLog[index].P
	for i := 0; i < len(pSet); i++ {
		found := false
		for j := 0; j < len(e); j++ {
			if e[j].Id == pSet[i].Id && e[j].Fit == pSet[i].Fit {
				found = true
				break
			}
		}
		if found == false {
			return false
		}
	}
	return true
}

/*
	extract the E sets from the messages
*/

func (in *Instance) getESetFromESets(responses []*proto.GenericConsensus) []*proto.GenericConsensusValue {
	eSet := make([]*proto.GenericConsensusValue, 0)
	for i := 0; i < len(responses); i++ {
		eSet = in.setUnionProtoValues(eSet, responses[i].E)
		eSet = in.removeEmptyValues(eSet)
	}
	return eSet
}

/*
	returns the best proposal in the array: best proposal is the proposal with the highest priority.
	Note that priorities have the form pi.proposer.recorder such that in case of equal pi, the higher node number wins
*/

func (in *Instance) getBestProposal(e []*proto.GenericConsensusValue) *proto.GenericConsensusValue {
	var hiValue *proto.GenericConsensusValue

	hiValue = e[0]

	for i := 1; i < len(e); i++ {
		if in.hasHigherPriority(hiValue, e[i]) {
			hiValue = e[i]
		}
	}

	return hiValue
}

/*
	checks if value2 has a higher priority than value1
*/

func (in *Instance) hasHigherPriority(value1 *proto.GenericConsensusValue, value2 *proto.GenericConsensusValue) bool {
	fit1 := value1.Fit
	fit2 := value2.Fit

	if in.debugOn {
		f1 := strings.Split(fit1, ".")
		f2 := strings.Split(fit2, ".")
		if len(f1) < 3 || len(f2) < 3 {
			panic("error in priorities:value 1: " + value1.Id + "-" + value1.Fit + ", value 2: " + value2.Id + "-" + value2.Fit)
		}
	}
	pi_1, _ := strconv.Atoi(strings.Split(fit1, ".")[0])
	pi_2, _ := strconv.Atoi(strings.Split(fit2, ".")[0])

	prop_1, _ := strconv.Atoi(strings.Split(fit1, ".")[1])
	prop_2, _ := strconv.Atoi(strings.Split(fit2, ".")[1])

	rec_1, _ := strconv.Atoi(strings.Split(fit1, ".")[2])
	rec_2, _ := strconv.Atoi(strings.Split(fit2, ".")[2])

	if pi_2 > pi_1 {
		return true
	} else if pi_1 > pi_2 {
		return false
	} else if pi_2 == pi_1 {
		if prop_2 > prop_1 {
			return true
		} else if prop_2 < prop_1 {
			return false
		} else if prop_2 == prop_1 {
			if rec_2 > rec_1 {
				return true
			} else if rec_2 < rec_1 {
				return false
			} else if rec_2 == rec_1 {
				panic("Found equal priorities: " + fit1 + " and " + fit2)
			}
		}
	}

	return false

}

/*
	upon receiving majority gather messages, increase s, set P to be best(c) and reset the variables
*/

func (in *Instance) proposerHandleGatherCQQuorum(consensusMessage *proto.GenericConsensus) {
	in.proposerReplicatedLog[consensusMessage.Index].S = in.proposerReplicatedLog[consensusMessage.Index].S + 1
	in.proposerReplicatedLog[consensusMessage.Index].P = make([]*proto.GenericConsensusValue, 1)
	in.proposerReplicatedLog[consensusMessage.Index].P[0] = &proto.GenericConsensusValue{
		Id:  in.getBestProposal(in.proposerReplicatedLog[consensusMessage.Index].C).Id,
		Fit: "",
	}
	in.proposerReplicatedLog[consensusMessage.Index].E = []*proto.GenericConsensusValue{}
	in.proposerReplicatedLog[consensusMessage.Index].C = []*proto.GenericConsensusValue{}

	in.proposerReplicatedLog[consensusMessage.Index].proposeResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].spreadEResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].spreadCGatherEResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].gatherCResponses = []*proto.GenericConsensus{}
}

/*
	Catch-up when received a message with a higher time stamp from the recorder
*/

func (in *Instance) consensusCatchUp(consensusMessage *proto.GenericConsensus) {
	in.proposerReplicatedLog[consensusMessage.Index].S = consensusMessage.S
	in.proposerReplicatedLog[consensusMessage.Index].P = consensusMessage.P
	in.proposerReplicatedLog[consensusMessage.Index].E = in.removeEmptyValues(consensusMessage.E)
	in.proposerReplicatedLog[consensusMessage.Index].C = in.removeEmptyValues(consensusMessage.C)

	in.proposerReplicatedLog[consensusMessage.Index].proposeResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].spreadEResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].spreadCGatherEResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].gatherCResponses = []*proto.GenericConsensus{}

	in.proposerSendMessage(consensusMessage.Index)

}
