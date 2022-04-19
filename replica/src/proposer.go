package raxos

import (
	"math"
	"raxos/proto"
	"strconv"
	"strings"
)

/*
	Proposer: handler for proposer consensus messages
*/

func (in *Instance) handleProposerConsensusMessage(consensusMessage *proto.GenericConsensus) {

	// case 2: a propose reply message from a recorder

	if consensusMessage.M == in.proposeMessage && in.proposerReplicatedLog[consensusMessage.Index].S == consensusMessage.S && in.proposerReplicatedLog[consensusMessage.Index].S%4 == 1 {
		in.proposerReplicatedLog[consensusMessage.Index].E = in.setUnionProtoValue(in.proposerReplicatedLog[consensusMessage.Index].E, consensusMessage.P)
		in.debug("Received propose response from recorder has id = "+consensusMessage.P.Id+" and fit = "+consensusMessage.P.Fit, 0)
		in.proposerReplicatedLog[consensusMessage.Index].E = in.removeEmptyValues(in.proposerReplicatedLog[consensusMessage.Index].E)
		in.proposerReplicatedLog[consensusMessage.Index].proposeResponses = append(in.proposerReplicatedLog[consensusMessage.Index].proposeResponses, consensusMessage)

		// case 2.1: on receiving quorum of propose replies
		if int64(len(in.proposerReplicatedLog[consensusMessage.Index].proposeResponses)) == in.numReplicas/2+1 {

			majorityValue := in.getProposalWithMajorityHi(in.proposerReplicatedLog[consensusMessage.Index].proposeResponses)

			// case 2.1.1 a majority of the proposals responses have Hi and I am not decided yet
			if majorityValue.Id != "" && majorityValue.Fit != "" && in.proposerReplicatedLog[consensusMessage.Index].decided == false {
				in.proposerReceivedMajorityProposeWithHi(consensusMessage, majorityValue)
				return
			}

			// case 2.1.2 no proposal with majority Hi
			if majorityValue.Id == "" && majorityValue.Fit == "" {
				in.proposerReplicatedLog[consensusMessage.Index].S = in.proposerReplicatedLog[consensusMessage.Index].S + 1
				in.proposerSendSpreadE(consensusMessage.Index)
				return
			}

		}
		return
	}

	// case 3: a spreadE response
	if consensusMessage.M == in.spreadEMessage && in.proposerReplicatedLog[consensusMessage.Index].S == consensusMessage.S && in.proposerReplicatedLog[consensusMessage.Index].S%4 == 2 {
		in.proposerReplicatedLog[consensusMessage.Index].spreadEResponses = append(in.proposerReplicatedLog[consensusMessage.Index].spreadEResponses, consensusMessage)

		// case 3.1 upon receiving a f+1 spreadE message
		if int64(len(in.proposerReplicatedLog[consensusMessage.Index].spreadEResponses)) == in.numReplicas/2+1 {
			EdashSet := in.getESetfromSpreadEResponses(in.proposerReplicatedLog[consensusMessage.Index].spreadEResponses) // EdashSet is a 2D array that contains the E' set extracted from each spreadE reponse
			// case 3.1.1 all E' sets contain Hi for the same value
			var hiValue *proto.GenericConsensusValue
			hiValue = in.getProposalWithMajoriyHiInAllESets(EdashSet)
			if hiValue.Id != "" && hiValue.Fit != "" && in.proposerReplicatedLog[consensusMessage.Index].decided == false {
				in.proposerReceivedMajorityProposeWithHi(consensusMessage, hiValue)
				return
			}
			// case 3.1.2 no hash has f+1 Hi
			if hiValue.Id == "" || hiValue.Fit == "" {
				in.proposerReplicatedLog[consensusMessage.Index].C = in.proposerReplicatedLog[consensusMessage.Index].E
				in.proposerReplicatedLog[consensusMessage.Index].C = in.removeEmptyValues(in.proposerReplicatedLog[consensusMessage.Index].C)
				in.proposerReplicatedLog[consensusMessage.Index].S = in.proposerReplicatedLog[consensusMessage.Index].S + 1
				in.proposerSendSpreadCGatherE(consensusMessage.Index)
				return
			}

		}
		return

	}

	// case 4: a SpreadCGatherE response
	if consensusMessage.M == in.spreadCgatherEMessage && in.proposerReplicatedLog[consensusMessage.Index].S == consensusMessage.S && in.proposerReplicatedLog[consensusMessage.Index].S%4 == 3 {
		// copy the E set
		in.proposerReplicatedLog[consensusMessage.Index].E = in.setUnionProtoValues(in.proposerReplicatedLog[consensusMessage.Index].E, consensusMessage.E)
		in.proposerReplicatedLog[consensusMessage.Index].E = in.removeEmptyValues(in.proposerReplicatedLog[consensusMessage.Index].E)
		in.proposerReplicatedLog[consensusMessage.Index].spreadCGatherEResponses = append(in.proposerReplicatedLog[consensusMessage.Index].spreadCGatherEResponses, consensusMessage)

		// case 4.1: upon receiving a majority SpreadCGatherE responses

		if len(in.proposerReplicatedLog[consensusMessage.Index].spreadCGatherEResponses) == int(in.numReplicas)/2+1 {
			in.proposerReplicatedLog[consensusMessage.Index].S = in.proposerReplicatedLog[consensusMessage.Index].S + 1
			in.proposerReplicatedLog[consensusMessage.Index].U = in.proposerReplicatedLog[consensusMessage.Index].C

			in.proposerReplicatedLog[consensusMessage.Index].U = in.removeEmptyValues(in.proposerReplicatedLog[consensusMessage.Index].U)
			in.proposerReplicatedLog[consensusMessage.Index].E = in.removeEmptyValues(in.proposerReplicatedLog[consensusMessage.Index].E)
			E_Best := in.getBestProposal(in.proposerReplicatedLog[consensusMessage.Index].E)
			//C_Best := in.getBestProposal(in.proposerReplicatedLog[consensusMessage.Index].C)
			U_Best := in.getBestProposal(in.proposerReplicatedLog[consensusMessage.Index].U)

			// case 4.1.1 If U.getTheBest() == E.getTheBest()
			if U_Best.Id != "" && U_Best.Fit != "" && E_Best.Id != "" && E_Best.Fit != "" {
				if U_Best.Id == E_Best.Id && U_Best.Fit == E_Best.Fit {
					in.debug("decided "+strconv.Itoa(int(consensusMessage.Index))+" in the SpreadCGatherE phase, step is "+strconv.Itoa(int(in.proposerReplicatedLog[consensusMessage.Index].S))+" size of E set is "+strconv.Itoa(len(in.proposerReplicatedLog[consensusMessage.Index].E)), 3)
					in.proposerReceivedMajorityProposeWithHi(consensusMessage, U_Best)
					return
				}
			}

			in.proposerSendGatherC(consensusMessage.Index)
			return

		}
		return
	}

	// case 5: a GatherC response
	if consensusMessage.M == in.gatherCMessage && in.proposerReplicatedLog[consensusMessage.Index].S == consensusMessage.S && in.proposerReplicatedLog[consensusMessage.Index].S%4 == 0 {
		// update the C set
		in.proposerReplicatedLog[consensusMessage.Index].C = in.setUnionProtoValues(in.proposerReplicatedLog[consensusMessage.Index].C, consensusMessage.C)
		in.proposerReplicatedLog[consensusMessage.Index].C = in.removeEmptyValues(in.proposerReplicatedLog[consensusMessage.Index].C)
		in.proposerReplicatedLog[consensusMessage.Index].gatherCResponses = append(in.proposerReplicatedLog[consensusMessage.Index].gatherCResponses, consensusMessage)

		// case 5.1: upon receiving majority GatherC

		if len(in.proposerReplicatedLog[consensusMessage.Index].gatherCResponses) == int(in.numReplicas)/2+1 {
			in.proposerResetSlotAfterReceivingGather(consensusMessage)
			in.proposerSendPropose(consensusMessage.Index)
			return
		}
		return

	}

	// case 5: a catch-up message
	if consensusMessage.S > in.proposerReplicatedLog[consensusMessage.Index].S {
		in.consensusCatchUp(consensusMessage)
	}
}

/*
	Catch-up when received a message with a higher time stamp from the recorder
*/

func (in *Instance) consensusCatchUp(consensusMessage *proto.GenericConsensus) {
	in.proposerReplicatedLog[consensusMessage.Index].S = consensusMessage.S
	in.proposerReplicatedLog[consensusMessage.Index].P = consensusMessage.P

	in.proposerReplicatedLog[consensusMessage.Index].E = in.removeEmptyValues(consensusMessage.E)
	in.proposerReplicatedLog[consensusMessage.Index].C = in.removeEmptyValues(consensusMessage.C)

	in.proposerReplicatedLog[consensusMessage.Index].U = []*proto.GenericConsensusValue{}

	in.proposerReplicatedLog[consensusMessage.Index].proposeResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].spreadEResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].spreadCGatherEResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].gatherCResponses = []*proto.GenericConsensus{}

	if consensusMessage.S%4 == 1 {
		//send a propose message
		in.proposerSendPropose(consensusMessage.Index)
		return
	} else if consensusMessage.S%4 == 2 {
		//send a spreadE message
		in.proposerSendSpreadE(consensusMessage.Index)
		return
	} else if consensusMessage.S%4 == 3 {
		//send a spreadCGatherE message
		in.proposerSendSpreadCGatherE(consensusMessage.Index)
		return
	} else if consensusMessage.S%4 == 0 {
		//send a gatherC message
		in.proposerSendGatherC(consensusMessage.Index)
		return
	}

}

/*
	Upon receiving a quorum of propose responses, broadcast a spreadE message to all recorders
*/

func (in *Instance) proposerSendSpreadE(index int64) {
	for i := int64(0); i < in.numReplicas; i++ {
		consensusSpreadE := proto.GenericConsensus{
			Sender:      in.nodeName,
			Receiver:    i,
			Index:       index,
			M:           in.spreadEMessage,
			S:           in.proposerReplicatedLog[index].S,
			P:           nil,
			E:           in.removeEmptyValues(in.proposerReplicatedLog[index].E),
			C:           nil,
			D:           false,
			DS:          nil,
			PR:          -1,
			Destination: in.consensusMessageRecorderDestination,
		}

		rpcPair := RPCPair{
			Code: in.genericConsensusRpc,
			Obj:  &consensusSpreadE,
		}
		in.debug("sending a spreadE consensus message to recorder "+strconv.Itoa(int(i)), 1)

		in.sendMessage(i, rpcPair)
	}
}

/*
	Proposer: Called when the proposer receives f+1 propose responses, and each response corresponds to Hi priority for the same value
*/

func (in *Instance) proposerReceivedMajorityProposeWithHi(consensusMessage *proto.GenericConsensus, majorityValue *proto.GenericConsensusValue) {
	in.proposerReplicatedLog[consensusMessage.Index].S = math.MaxInt64
	in.proposerReplicatedLog[consensusMessage.Index].decided = true
	in.proposerReplicatedLog[consensusMessage.Index].decision = majorityValue
	in.proposerReplicatedLog[consensusMessage.Index].proposer = in.getProposer(majorityValue)

	// clean the memory
	in.proposerReplicatedLog[consensusMessage.Index].P = &proto.GenericConsensusValue{}
	in.proposerReplicatedLog[consensusMessage.Index].E = []*proto.GenericConsensusValue{}
	in.proposerReplicatedLog[consensusMessage.Index].C = []*proto.GenericConsensusValue{}
	in.proposerReplicatedLog[consensusMessage.Index].U = []*proto.GenericConsensusValue{}
	in.proposerReplicatedLog[consensusMessage.Index].proposeResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].spreadEResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].spreadCGatherEResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].gatherCResponses = []*proto.GenericConsensus{}

	in.delivered(consensusMessage.Index, majorityValue.Id, in.proposerReplicatedLog[consensusMessage.Index].proposer)

	// send a decide message

	for i := int64(0); i < in.numReplicas; i++ {

		consensusDecide := proto.GenericConsensus{
			Sender:      in.nodeName,
			Receiver:    i,
			Index:       consensusMessage.Index,
			M:           in.decideMessage,
			S:           in.proposerReplicatedLog[consensusMessage.Index].S, // not needed
			P:           nil,
			E:           nil,
			C:           nil,
			D:           true,
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
	Scan the propose replies set from recorders and assign the number of times the Hi priority appears for each proposal
*/

func (in *Instance) getProposalWithMajorityHi(e []*proto.GenericConsensus) *proto.GenericConsensusValue {
	var ValueCount map[string]int64 // maps each id to hi count
	ValueCount = make(map[string]int64)
	for i := 0; i < len(e); i++ {
		if strings.HasPrefix(e[i].P.Fit, strconv.FormatInt(in.Hi, 10)) {
			_, ok := ValueCount[e[i].P.Id+e[i].P.Fit] // note that this is a unique key
			if !ok {
				ValueCount[e[i].P.Id+e[i].P.Fit] = 1 // note that each proposal has a unique id
			} else {
				ValueCount[e[i].P.Id+e[i].P.Fit] = ValueCount[e[i].P.Id+e[i].P.Fit] + 1
			}
		}
	}
	// check whether there is a value who has f+1 Hi count

	for key, element := range ValueCount {
		if element >= in.numReplicas/2+1 {
			for i := 0; i < len(e); i++ {
				if e[i].P.Id+e[i].P.Fit == key {
					return e[i].P
				}
			}
		}
	}

	return &proto.GenericConsensusValue{
		Id:  "",
		Fit: "",
	} // there is no proposal with f+1 Hi

}

/*
	Proposer: mark the slot as decided and inform the upper layer
*/

func (in *Instance) recordProposerDecide(consensusMessage *proto.GenericConsensus) {
	in.proposerReplicatedLog[consensusMessage.Index].S = consensusMessage.S
	in.proposerReplicatedLog[consensusMessage.Index].decided = true
	in.proposerReplicatedLog[consensusMessage.Index].decision = consensusMessage.DS
	in.proposerReplicatedLog[consensusMessage.Index].proposer = consensusMessage.PR

	in.proposerReplicatedLog[consensusMessage.Index].P = &proto.GenericConsensusValue{}
	in.proposerReplicatedLog[consensusMessage.Index].E = []*proto.GenericConsensusValue{}
	in.proposerReplicatedLog[consensusMessage.Index].C = []*proto.GenericConsensusValue{}
	in.proposerReplicatedLog[consensusMessage.Index].U = []*proto.GenericConsensusValue{}
	in.proposerReplicatedLog[consensusMessage.Index].proposeResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].spreadEResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].spreadCGatherEResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].gatherCResponses = []*proto.GenericConsensus{}

	in.delivered(consensusMessage.Index, consensusMessage.DS.Id, consensusMessage.PR)

}

/*
	propose a new block for the slot index
*/

func (in *Instance) propose(index int64, hash string) {
	/*
		We use pipelining, so its possible that the proposer proposes messages corresponding to different instances without receiving the decisions for the prior instances
	*/
	in.proposerReplicatedLog = in.initializeSlot(in.proposerReplicatedLog, index) // create the slot if not already created
	in.proposerReplicatedLog[index].S = 1
	in.proposerReplicatedLog[index].P = &proto.GenericConsensusValue{
		Id:  hash,
		Fit: "",
	}

	// send the proposal message to all the recorders
	in.proposerSendPropose(index)
}

/*
	indication from the consensus layer that a value is decided
*/

func (in *Instance) delivered(index int64, hash string, proposer int64) {
	in.updateStateMachine()
	if int64(len(in.proposed)) > index && in.proposed[index] != "" && in.proposed[index] != hash {
		/*I previously proposed for this index, but somebody else has won the consensus*/
		in.proposed[index] = ""
		in.updateProposedIndex(hash)
		in.propose(in.proposedIndex, hash)
		in.debug("Re proposed because the decided value is not the one I proposed earlier ", 1)
	}
}

/*
	Proposer: invoked whenever a new slot is decided
	marks the slot as committed
	sends the response back to the client
    //todo how does pipeline affect this? when failures occur slot i will be completed before slot i-n, in that case should we implement a new failure detector and repurpose?
*/

func (in *Instance) updateStateMachine() {
	for len(in.proposerReplicatedLog) > int(in.committedIndex)+1 && in.proposerReplicatedLog[in.committedIndex+1].decided == true {

		decision := in.proposerReplicatedLog[in.committedIndex+1].decision.Id
		in.debug("Decided "+strconv.Itoa(int(in.committedIndex+1))+" with decision "+decision, 1)
		// if the decision is a sequence of hashes, then invoke the following for each hash: check if all the hashes exist and succeed only if everything exists
		blocks := make([]*proto.MessageBlock, 0)
		hashes := strings.Split(decision, ":")
		for i := 0; i < len(hashes); i++ {
			hashes[i] = strings.TrimSpace(hashes[i])
			in.debug("hash "+strconv.Itoa(i)+" is "+hashes[i], 0)
		}

		for i := 0; i < len(hashes); i++ {
			if len(hashes[i]) == 0 {
				continue
			}
			messageBlock, ok := in.messageStore.Get(hashes[i])
			if !ok {
				in.debug("Couldn't find hash "+hashes[i], 0)
				/*
					I haven't received the actual message block still, request it
				*/
				in.sendMessageBlockRequest(hashes[i])
				return
			}
			blocks = append(blocks, messageBlock)
		}
		// found all the blocks
		in.proposerReplicatedLog[in.committedIndex+1].committed = true
		in.debug("Committed "+strconv.Itoa(int(in.committedIndex+1)), 2)
		in.committedIndex++
		in.numInflightRequests--

		for i := 0; i < len(blocks); i++ {
			in.executeAndSendResponse(blocks[i])
		}
	}
}

/*
	Execute each command in the block, get the response array and send the responses back to the client if I am the proposer who created this block
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
	Upon receiving a consensus request, the leader node will propose it to the consensus layer if the number of inflight requests is less than the pipeline length
*/

func (in *Instance) handleConsensusRequest(request *proto.ConsensusRequest) {
	if len([]rune(request.Hash)) > 0 {
		//if in.nodeName == in.getDeterministicLeader1() {
		if true { // for testing purpose
			if true { //todo in.numInflightRequests <= in.pipelineLength { // this should be changed to not drop requests
				in.updateProposedIndex(request.Hash)
				in.propose(in.proposedIndex, request.Hash)
				in.numInflightRequests++
				in.debug("Proposed "+request.Hash+" to "+strconv.Itoa(int(in.proposedIndex)), 0)
			} else {
				in.debug("Proposed failed due to inflight requests"+request.Hash+" to "+strconv.Itoa(int(in.proposedIndex)), 1)
			}
		}
	}
}

/*
	Proposer: Set proposed index = max(proposed_index, committed index)+1
	add "" entries upto proposed index
    set proposed[index]=hash
*/

func (in *Instance) updateProposedIndex(hash string) {
	maxIndex := in.proposedIndex
	if in.committedIndex > maxIndex {
		maxIndex = in.committedIndex
	}
	in.proposedIndex = maxIndex + 1
	for i := int64(len(in.proposed)); i < in.proposedIndex+1; i++ {
		in.proposed = append(in.proposed, "")
	}
	in.proposed[in.proposedIndex] = hash
	in.debug("Added "+hash+" to proposed array position "+strconv.Itoa(int(in.proposedIndex)), 1)
}

/*
	Proposer: Returns the proposer who proposed this value, extracted using the fitness
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
	A util function to create an E multi set which is a 2d array, each element of E is the E set of the response[i]
*/

func (in *Instance) getESetfromSpreadEResponses(responses []*proto.GenericConsensus) [][]*proto.GenericConsensusValue {
	var EMultiSet [][]*proto.GenericConsensusValue
	EMultiSet = make([][]*proto.GenericConsensusValue, len(responses))
	for i := 0; i < len(responses); i++ {
		EMultiSet[i] = make([]*proto.GenericConsensusValue, 0)
		EMultiSet[i] = in.setUnionProtoValues(EMultiSet[i], responses[i].E)
	}
	return EMultiSet
}

/*
	Proposer: Scans all the E sets, and checks if there is a proposal which appears in all the E sets with priority Hi
*/

func (in *Instance) getProposalWithMajoriyHiInAllESets(set [][]*proto.GenericConsensusValue) *proto.GenericConsensusValue {
	hiValue := &proto.GenericConsensusValue{}

	var HiCounts map[string]int64
	HiCounts = make(map[string]int64)

	for i := 0; i < len(set); i++ {
		seenHashes := make(map[string]bool) // we need to make sure that we count only one occurrence of the same proposal in each set E (if a given proposal appears more than once in an E set, we only count once)
		for j := 0; j < len(set[i]); j++ {
			_, ok := seenHashes[set[i][j].Id+set[i][j].Fit]
			if ok {
				continue
			} else {
				seenHashes[set[i][j].Id+set[i][j].Fit] = true
				if strings.HasPrefix(set[i][j].Fit, strconv.FormatInt(in.Hi, 10)) {
					_, ok = HiCounts[set[i][j].Id+set[i][j].Fit]
					if ok {
						HiCounts[set[i][j].Id+set[i][j].Fit] = HiCounts[set[i][j].Id+set[i][j].Fit] + 1
					} else {
						HiCounts[set[i][j].Id+set[i][j].Fit] = 1
					}
				}

			}
		}
	}

	for key, element := range HiCounts {
		if element >= in.numReplicas/2+1 {

			for i := 0; i < len(set); i++ {
				for j := 0; j < len(set[i]); j++ {
					if set[i][j].Id+set[i][j].Fit == key {
						return set[i][j]
					}
				}
			}

		}
	}

	return hiValue
}

/*
	Proposer: Send a SpreadCGatherE Message to all recorders
*/

func (in *Instance) proposerSendSpreadCGatherE(index int64) {
	for i := int64(0); i < in.numReplicas; i++ {
		consensusSpreadCGatherE := proto.GenericConsensus{
			Sender:      in.nodeName,
			Receiver:    i,
			Index:       index,
			M:           in.spreadCgatherEMessage,
			S:           in.proposerReplicatedLog[index].S,
			P:           nil,
			E:           nil,
			C:           in.removeEmptyValues(in.proposerReplicatedLog[index].C),
			D:           false,
			DS:          nil,
			PR:          -1,
			Destination: in.consensusMessageRecorderDestination,
		}

		rpcPair := RPCPair{
			Code: in.genericConsensusRpc,
			Obj:  &consensusSpreadCGatherE,
		}
		in.debug("sending a spreadCGatherE consensus message to "+strconv.Itoa(int(i)), 1)

		in.sendMessage(i, rpcPair)
	}
}

/*
	Returns the best proposal in the array: best proposal is the proposal with the highest priority.
	Note that priorities have the form pi.node such that in case of equal pi, the higher node number wins
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
			//os.Exit(255)
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
	Send a GatherC message to all recorders
*/

func (in *Instance) proposerSendGatherC(index int64) {
	for i := int64(0); i < in.numReplicas; i++ {
		consensusGatherC := proto.GenericConsensus{
			Sender:      in.nodeName,
			Receiver:    i,
			Index:       index,
			M:           in.gatherCMessage,
			S:           in.proposerReplicatedLog[index].S,
			P:           nil,
			E:           nil,
			C:           nil,
			D:           false,
			DS:          nil,
			PR:          -1,
			Destination: in.consensusMessageRecorderDestination,
		}

		rpcPair := RPCPair{
			Code: in.genericConsensusRpc,
			Obj:  &consensusGatherC,
		}
		in.debug("sending a gatherC consensus message to "+strconv.Itoa(int(i)), 1)

		in.sendMessage(i, rpcPair)
	}
}

/*
	Broadcast slow path propose message to all recorders: what is slow path? we do not make any distinction in the recorder side
*/

func (in *Instance) proposerSendPropose(index int64) {
	for i := int64(0); i < in.numReplicas; i++ {

		consensusPropose := proto.GenericConsensus{
			Sender:      in.nodeName,
			Receiver:    i,
			Index:       index,
			M:           in.proposeMessage,
			S:           in.proposerReplicatedLog[index].S,
			P:           in.proposerReplicatedLog[index].P,
			E:           nil,
			C:           nil,
			D:           false,
			DS:          nil,
			PR:          -1,
			Destination: in.consensusMessageRecorderDestination,
		}

		rpcPair := RPCPair{
			Code: in.genericConsensusRpc,
			Obj:  &consensusPropose,
		}
		in.debug("sending a consensus propose message to "+strconv.Itoa(int(i)), 1)

		in.sendMessage(i, rpcPair)
	}
}

/*
	Upon receiving majority gather messages, increase s and reset the variables
*/

func (in *Instance) proposerResetSlotAfterReceivingGather(consensusMessage *proto.GenericConsensus) {
	in.proposerReplicatedLog[consensusMessage.Index].S = in.proposerReplicatedLog[consensusMessage.Index].S + 1
	in.proposerReplicatedLog[consensusMessage.Index].P = &proto.GenericConsensusValue{
		Id:  in.getBestProposal(in.proposerReplicatedLog[consensusMessage.Index].C).Id,
		Fit: "",
	}
	in.proposerReplicatedLog[consensusMessage.Index].E = []*proto.GenericConsensusValue{}
	in.proposerReplicatedLog[consensusMessage.Index].C = []*proto.GenericConsensusValue{}
	in.proposerReplicatedLog[consensusMessage.Index].U = []*proto.GenericConsensusValue{}

	in.proposerReplicatedLog[consensusMessage.Index].proposeResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].spreadEResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].spreadCGatherEResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].gatherCResponses = []*proto.GenericConsensus{}
}
