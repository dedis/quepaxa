package raxos

import (
	"math"
	"raxos/proto"
	"strconv"
	"strings"
)

/*
	Proposer: This is a message from the recorder to the proposer
*/

func (in *Instance) handleProposerConsensusMessage(consensusMessage *proto.GenericConsensus) {

	// case 1: a decide message from a recorder. This is not handled here, but in the common.go since it is common to both the proposer and the recorder

	// case 2: a propose reply message from a recorder

	if consensusMessage.M == in.proposeMessage && in.proposerReplicatedLog[consensusMessage.Index].S == consensusMessage.S {
		in.proposerReplicatedLog[consensusMessage.Index].E = append(in.proposerReplicatedLog[consensusMessage.Index].E, Value{
			id:  consensusMessage.P.Id,
			fit: consensusMessage.P.Fit,
		})

		in.proposerReplicatedLog[consensusMessage.Index].proposeResponses = append(in.proposerReplicatedLog[consensusMessage.Index].proposeResponses, consensusMessage)

		// case 2.1: on receiving quorum of propose replies
		if int64(len(in.proposerReplicatedLog[consensusMessage.Index].proposeResponses)) == in.numReplicas/2+1 {

			majorityValue := in.getProposalWithMajorityHi(in.proposerReplicatedLog[consensusMessage.Index].proposeResponses)

			// case 2.1.1 a majority of the proposals responses have Hi and I am not decided yet
			if majorityValue.id != "" && majorityValue.fit != "" && in.proposerReplicatedLog[consensusMessage.Index].decided == false {
				in.proposerReceivedMajorityProposeWithHi(consensusMessage, majorityValue)
				return
			}

			// case 2.1.2 no proposal with majority Hi
			if majorityValue.id == "" && majorityValue.fit == "" {
				in.proposerReplicatedLog[consensusMessage.Index].S = in.proposerReplicatedLog[consensusMessage.Index].S + 1
				in.proposerSendSpreadE(consensusMessage.Index)
				return
			}

		}
		return
	}

	// case 3: a spreadE response
	if consensusMessage.M == in.spreadEMessage && in.proposerReplicatedLog[consensusMessage.Index].S == consensusMessage.S {
		in.proposerReplicatedLog[consensusMessage.Index].spreadEResponses = append(in.proposerReplicatedLog[consensusMessage.Index].spreadEResponses, consensusMessage)

		// case 3.1 upon receiving a f+1 spreadE message
		if int64(len(in.proposerReplicatedLog[consensusMessage.Index].spreadEResponses)) == in.numReplicas/2+1 {
			EdashSet := in.getESetfromSpreadEResponses(in.proposerReplicatedLog[consensusMessage.Index].spreadEResponses) // EdashSet is a 2D array that contains the E' set extracted from each spreadE reponse
			// case 3.1.1 all E' sets contain Hi for the same value
			var hiValue Value
			hiValue = in.getProposalWithMajoriyHiInAllESets(EdashSet)
			if hiValue.id != "" && hiValue.fit != "" && in.proposerReplicatedLog[consensusMessage.Index].decided == false {
				in.proposerReceivedMajorityProposeWithHi(consensusMessage, hiValue)
				return
			}
			// case 3.1.2 no hash has f+1 Hi
			if hiValue.id == "" || hiValue.fit == "" {
				in.proposerReplicatedLog[consensusMessage.Index].C = in.proposerReplicatedLog[consensusMessage.Index].E
				in.proposerReplicatedLog[consensusMessage.Index].S = in.proposerReplicatedLog[consensusMessage.Index].S + 1
				in.proposerSendSpreadCGatherE(consensusMessage.Index)
				return
			}

		}
		return

	}

	// case 4: a SpreadCGatherE response
	if consensusMessage.M == in.spreadCgatherEMessage && in.proposerReplicatedLog[consensusMessage.Index].S == consensusMessage.S {
		// copy the E set
		for i := 0; i < len(consensusMessage.E); i++ {
			in.proposerReplicatedLog[consensusMessage.Index].E = append(in.proposerReplicatedLog[consensusMessage.Index].E, Value{
				id:  consensusMessage.E[i].Id,
				fit: consensusMessage.E[i].Fit,
			})
		}
		in.proposerReplicatedLog[consensusMessage.Index].spreadCGatherEResponses = append(in.proposerReplicatedLog[consensusMessage.Index].spreadCGatherEResponses, consensusMessage)

		// case 4.1: upon receiving a majority SpreadCGatherE responses

		if len(in.proposerReplicatedLog[consensusMessage.Index].spreadCGatherEResponses) == int(in.numReplicas)/2+1 {
			in.proposerReplicatedLog[consensusMessage.Index].S = in.proposerReplicatedLog[consensusMessage.Index].S + 1
			in.proposerReplicatedLog[consensusMessage.Index].U = in.proposerReplicatedLog[consensusMessage.Index].C
			E_Best := in.getBestProposal(in.proposerReplicatedLog[consensusMessage.Index].E)
			//C_Best := in.getBestProposal(in.proposerReplicatedLog[consensusMessage.Index].C)
			U_Best := in.getBestProposal(in.proposerReplicatedLog[consensusMessage.Index].U)

			// case 4.1.1 If U.getTheBest() == E.getTheBest()
			if U_Best.id != "" && U_Best.fit != "" && E_Best.id != "" && E_Best.fit != "" {
				if U_Best.id == E_Best.id && U_Best.fit == E_Best.fit {
					in.proposerReceivedMajorityProposeWithHi(consensusMessage, U_Best) //todo C-best or U-best?
					return
				}
			}

			in.proposerSendGatherC(consensusMessage.Index)
			return

		}
		return
	}

	// case 5: a GatherC response
	if consensusMessage.M == in.gatherCMessage && in.proposerReplicatedLog[consensusMessage.Index].S == consensusMessage.S {
		// update the C set
		for i := 0; i < len(consensusMessage.C); i++ {
			in.proposerReplicatedLog[consensusMessage.Index].C = append(in.proposerReplicatedLog[consensusMessage.Index].C, Value{
				id:  consensusMessage.C[i].Id,
				fit: consensusMessage.C[i].Fit,
			})
		}
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
	Proposer: catch-up when received a message with a higher time stamp from the recorder
*/

func (in *Instance) consensusCatchUp(consensusMessage *proto.GenericConsensus) {
	in.proposerReplicatedLog[consensusMessage.Index].S = consensusMessage.S
	in.proposerReplicatedLog[consensusMessage.Index].P = Value{
		id:  consensusMessage.P.Id,
		fit: consensusMessage.P.Fit,
	}
	in.proposerReplicatedLog[consensusMessage.Index].E = []Value{}
	in.proposerReplicatedLog[consensusMessage.Index].C = []Value{}
	in.proposerReplicatedLog[consensusMessage.Index].U = []Value{}

	in.proposerReplicatedLog[consensusMessage.Index].proposeResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].spreadEResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].spreadCGatherEResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].gatherCResponses = []*proto.GenericConsensus{}

	in.proposerReplicatedLog[consensusMessage.Index].E = in.proposerConvertToValueArray(consensusMessage.E)
	in.proposerReplicatedLog[consensusMessage.Index].C = in.proposerConvertToValueArray(consensusMessage.C)

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
	Proposer: Upon receiving a quarum of propose responses, broadcast a spreadE message to all recorders
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
			E:           in.getGenericConsensusValueArray(in.proposerReplicatedLog[index].E),
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
		in.debug("sending a spreadE consensus message to "+strconv.Itoa(int(i)), 1)

		in.sendMessage(i, rpcPair)
	}
}

/*
	Proposer: Called when the proposer receives f+1 propose responses, and each response correspond to Hi priority for the same value
*/

func (in *Instance) proposerReceivedMajorityProposeWithHi(consensusMessage *proto.GenericConsensus, majorityValue Value) {
	in.proposerReplicatedLog[consensusMessage.Index].S = math.MaxInt64
	in.proposerReplicatedLog[consensusMessage.Index].decided = true
	in.proposerReplicatedLog[consensusMessage.Index].decision = Value{
		id:  majorityValue.id,
		fit: majorityValue.fit,
	}
	in.proposerReplicatedLog[consensusMessage.Index].proposer = in.getProposer(majorityValue)

	in.proposerReplicatedLog[consensusMessage.Index].P = Value{}
	in.proposerReplicatedLog[consensusMessage.Index].E = []Value{}
	in.proposerReplicatedLog[consensusMessage.Index].C = []Value{}
	in.proposerReplicatedLog[consensusMessage.Index].U = []Value{}
	in.proposerReplicatedLog[consensusMessage.Index].proposeResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].spreadEResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].spreadCGatherEResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].gatherCResponses = []*proto.GenericConsensus{}

	in.delivered(consensusMessage.Index, majorityValue.id, in.proposerReplicatedLog[consensusMessage.Index].proposer)

	// send a decide message

	for i := int64(0); i < in.numReplicas; i++ {

		consensusDecide := proto.GenericConsensus{
			Sender:   in.nodeName,
			Receiver: i,
			Index:    consensusMessage.Index,
			M:        in.decideMessage,
			S:        in.proposerReplicatedLog[consensusMessage.Index].S,
			P:        nil,
			E:        nil,
			C:        nil,
			D:        true,
			DS: &proto.GenericConsensusValue{
				Id:  in.proposerReplicatedLog[consensusMessage.Index].decision.id,
				Fit: in.proposerReplicatedLog[consensusMessage.Index].decision.fit,
			},
			PR:          in.proposerReplicatedLog[consensusMessage.Index].proposer,
			Destination: in.consensusMessageCommonDestination,
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
	Proposer: Scan the propose replies set and assign the number of times the Hi priority appears for each proposal
*/

func (in *Instance) getProposalWithMajorityHi(e []*proto.GenericConsensus) Value {
	var ValueCount map[string]int64 // maps each id to hi count
	ValueCount = make(map[string]int64)
	for i := 0; i < len(e); i++ {
		if strings.HasPrefix(e[i].P.Fit, strconv.FormatInt(in.Hi, 10)) {
			_, ok := ValueCount[e[i].P.Id+e[i].P.Fit]
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
					return Value{
						id:  e[i].P.Id,
						fit: e[i].P.Fit,
					}
				}
			}
		}
	}

	return Value{
		id:  "",
		fit: "",
	} // there is no proposal with f+1 Hi

}

/*
	Proposer: mark the slot as decided and inform the upper layer
*/

func (in *Instance) recordProposerDecide(consensusMessage *proto.GenericConsensus) {
	in.proposerReplicatedLog[consensusMessage.Index].S = consensusMessage.S
	in.proposerReplicatedLog[consensusMessage.Index].decided = true
	in.proposerReplicatedLog[consensusMessage.Index].decision = Value{
		id:  consensusMessage.DS.Id,
		fit: consensusMessage.DS.Fit,
	}
	in.proposerReplicatedLog[consensusMessage.Index].proposer = consensusMessage.PR

	in.proposerReplicatedLog[consensusMessage.Index].P = Value{}
	in.proposerReplicatedLog[consensusMessage.Index].E = []Value{}
	in.proposerReplicatedLog[consensusMessage.Index].C = []Value{}
	in.proposerReplicatedLog[consensusMessage.Index].U = []Value{}
	in.proposerReplicatedLog[consensusMessage.Index].proposeResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].spreadEResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].spreadCGatherEResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].gatherCResponses = []*proto.GenericConsensus{}

	in.delivered(consensusMessage.Index, consensusMessage.DS.Id, consensusMessage.PR)
}

/*
	Proposer: propose a new block for the slot index
*/

func (in *Instance) propose(index int64, hash string) {
	/*
		We use pipelining, so its possible that the proposer proposes messages corresponding to different instances without receiving the decisions for the prior instances
	*/
	in.proposerReplicatedLog = in.initializeSlot(in.proposerReplicatedLog, index) // create the slot if not already created
	in.proposerReplicatedLog[index].S = 1
	in.proposerReplicatedLog[index].P = Value{
		id:  hash,
		fit: "",
	}

	// send the proposal message to all the recorders
	in.proposerSendPropose(index)
}

/*
	Proposer: indication from the consensus layer that a value is decided
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
		decision := in.proposerReplicatedLog[in.committedIndex+1].decision.id
		// todo: if the decision is a sequence of hashes, then invoke the following for each hash: check if all the hashes exist and succeed only if everything exists
		messageBlock, ok := in.messageStore.Get(decision)
		if !ok {
			/*
				I haven't received the actual message block still, request it
			*/
			in.sendMessageBlockRequest(decision)
			return
		}

		in.proposerReplicatedLog[in.committedIndex+1].committed = true
		in.debug("Committed "+strconv.Itoa(int(in.committedIndex+1)), 1)
		in.committedIndex++
		in.numInflightRequests--
		in.executeAndSendResponse(messageBlock)
	}
}

/*
	Proposer: Execute each command in the block, get the response array and send the responses back to the client if I am the proposer who created this block
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
	Proposer: Send a consensus request to the leader / set of leaders. Upon receiving this, the leader node will eventually propose a value for a slot, for this hash
*/

func (in *Instance) sendConsensusRequest(hash string) {
	leader := in.getDeterministicLeader1() //todo change this when adding the improvements
	consensusRequest := proto.ConsensusRequest{
		Sender:   in.nodeName,
		Receiver: leader,
		Hash:     hash,
	}
	rpcPair := RPCPair{
		Code: in.consensusRequestRpc,
		Obj:  &consensusRequest,
	}
	in.debug("sending a consensus request to "+strconv.Itoa(int(leader))+" for the hash "+hash, 1)

	in.sendMessage(leader, rpcPair)
}

/*
	Proposer: Upon receiving a consensus request, the leader node will propose it to the consensus layer if the number of inflight requests is less than the pipeline length
*/

func (in *Instance) handleConsensusRequest(request *proto.ConsensusRequest) {
	if in.nodeName == in.getDeterministicLeader1() {
		if true { //in.numInflightRequests <= in.pipelineLength { // this should be changed to not drop requests
			in.updateProposedIndex(request.Hash)
			in.propose(in.proposedIndex, request.Hash)
			in.numInflightRequests++
			in.debug("Proposed "+request.Hash+" to "+strconv.Itoa(int(in.proposedIndex)), 1)
		} else {
			in.debug("Proposed failed due to inflight requests"+request.Hash+" to "+strconv.Itoa(int(in.proposedIndex)), 1)
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
	Proposer: Returns the proposer who proposed this value
*/

func (in *Instance) getProposer(value Value) int64 {
	priority := value.fit
	proposer, err := strconv.Atoi(strings.Split(priority, ".")[1])
	if err == nil {
		return int64(proposer)
	} else {
		return -1
	}
}

/*
	Proposer: Create an E multi set which is a 2d array, each element of E is the E set of the response[i]
*/

func (in *Instance) getESetfromSpreadEResponses(responses []*proto.GenericConsensus) [][]Value {
	var EMultiSet [][]Value
	EMultiSet = make([][]Value, len(responses))
	for i := 0; i < len(responses); i++ {
		EMultiSet[i] = make([]Value, len(responses[i].E))
		for j := 0; j < len(responses[i].E); j++ {
			EMultiSet[i] = append(EMultiSet[i], Value{
				id:  responses[i].E[j].Id,
				fit: responses[i].E[j].Fit,
			})
		}
	}
	return EMultiSet
}

/*
	Proposer: Scans all the spreadE E sets, and checks if there is a proposal which appears in all the E sets with priority Hi
*/

func (in *Instance) getProposalWithMajoriyHiInAllESets(set [][]Value) Value {
	hiValue := Value{
		id:  "",
		fit: "",
	}

	var HiCounts map[string]int64
	HiCounts = make(map[string]int64)

	for i := 0; i < len(set); i++ {
		seenHashes := make(map[string]bool) // we need to make sure that we count only one occurrence of the same proposal in each set E (if a given proposal appears more than once in an E set, we only count once)
		for j := 0; j < len(set[i]); j++ {
			_, ok := seenHashes[set[i][j].id+set[i][j].fit]
			if ok {
				continue
			} else {
				seenHashes[set[i][j].id+set[i][j].fit] = true
				if strings.HasPrefix(set[i][j].fit, strconv.FormatInt(in.Hi, 10)) {
					_, ok = HiCounts[set[i][j].id+set[i][j].fit]
					if ok {
						HiCounts[set[i][j].id+set[i][j].fit] = HiCounts[set[i][j].id+set[i][j].fit] + 1
					} else {
						HiCounts[set[i][j].id+set[i][j].fit] = 1
					}
				}

			}
		}
	}

	for key, element := range HiCounts {
		if element >= in.numReplicas/2+1 {

			for i := 0; i < len(set); i++ {
				for j := 0; j < len(set[i]); j++ {
					if set[i][j].id+set[i][j].fit == key {
						return Value{
							id:  set[i][j].id,
							fit: set[i][j].fit,
						}
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
			C:           in.getGenericConsensusValueArray(in.proposerReplicatedLog[index].C),
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
	Proposer: Returns the best proposal in the array: best proposal is the proposal with the highest priority.
	Note that priorities have the form pi.node such that in case of equal pi, the higher node number wins
*/

func (in *Instance) getBestProposal(e []Value) Value {
	hiValue := Value{
		id:  "",
		fit: "",
	}

	hiValue = e[0]

	for i := 1; i < len(e); i++ {
		if in.hasHigherPriority(hiValue, e[i]) {
			hiValue = e[i]
		}
	}

	return hiValue
}

/*
	Proposer: checks if value2 has a higher priority than value1
*/

func (in *Instance) hasHigherPriority(value1 Value, value2 Value) bool {
	fit1 := value1.fit
	fit2 := value2.fit

	pi_1, _ := strconv.Atoi(strings.Split(fit1, ".")[0])
	pi_2, _ := strconv.Atoi(strings.Split(fit2, ".")[0])

	node_1, _ := strconv.Atoi(strings.Split(fit1, ".")[1])
	node_2, _ := strconv.Atoi(strings.Split(fit2, ".")[1])

	if pi_2 > pi_1 {
		return true
	} else if pi_1 > pi_2 {
		return false
	} else if pi_2 == pi_1 {
		if node_2 > node_1 {
			return true
		} else {
			return false
		}
	}

	return false

}

/*
	Proposer: Send a GatherC message to all recorders
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
	Proposer: Broadcast slow path propose message to all recorders
*/

func (in *Instance) proposerSendPropose(index int64) {
	for i := int64(0); i < in.numReplicas; i++ {

		consensusPropose := proto.GenericConsensus{
			Sender:   in.nodeName,
			Receiver: i,
			Index:    index,
			M:        in.proposeMessage,
			S:        in.proposerReplicatedLog[index].S,
			P: &proto.GenericConsensusValue{
				Id:  in.proposerReplicatedLog[index].P.id,
				Fit: in.proposerReplicatedLog[index].P.fit,
			},
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
	Proposer: upon receiving majority gather messages, reset the variables
*/

func (in *Instance) proposerResetSlotAfterReceivingGather(consensusMessage *proto.GenericConsensus) {
	in.proposerReplicatedLog[consensusMessage.Index].S = in.proposerReplicatedLog[consensusMessage.Index].S + 1
	in.proposerReplicatedLog[consensusMessage.Index].P = Value{
		id:  in.getBestProposal(in.proposerReplicatedLog[consensusMessage.Index].C).id,
		fit: "",
	}
	in.proposerReplicatedLog[consensusMessage.Index].E = []Value{}
	in.proposerReplicatedLog[consensusMessage.Index].C = []Value{}
	in.proposerReplicatedLog[consensusMessage.Index].U = []Value{}

	in.proposerReplicatedLog[consensusMessage.Index].proposeResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].spreadEResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].spreadCGatherEResponses = []*proto.GenericConsensus{}
	in.proposerReplicatedLog[consensusMessage.Index].gatherCResponses = []*proto.GenericConsensus{}
}
