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

	// case 1: a decide message from a recorder

	if consensusMessage.M == in.decideMessage && in.proposerReplicatedLog[consensusMessage.Index].decided == false {
		in.recordProposerDecide(consensusMessage)
		return
	}

	// case 2: a propose reply message from a recorder

	if consensusMessage.M == in.proposeMessage && in.proposerReplicatedLog[consensusMessage.Index].S == consensusMessage.S {
		in.proposerReplicatedLog[consensusMessage.Index].E = append(in.proposerReplicatedLog[consensusMessage.Index].E, Value{
			id:  consensusMessage.P.Id,
			fit: consensusMessage.P.Fit,
		})

		// case 2.1: on receiving qaurum of replies
		if int64(len(in.proposerReplicatedLog[consensusMessage.Index].E)) >= in.numReplicas/2+1 {

			majorityValue := in.getProposalWithMajorityHi(in.proposerReplicatedLog[consensusMessage.Index].E)

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
		in.debug("sending a spreadE consensus message to " + strconv.Itoa(int(i)))

		in.sendMessage(i, rpcPair)
	}
}

/*
	Proposer: Called when the proposer receives f+1 propose responses, and each reponse correspond to Hi priority for the same value
*/

func (in *Instance) proposerReceivedMajorityProposeWithHi(consensusMessage *proto.GenericConsensus, majorityValue Value) {
	in.proposerReplicatedLog[consensusMessage.Index].S = math.MaxInt64
	in.proposerReplicatedLog[consensusMessage.Index].decided = true
	in.proposerReplicatedLog[consensusMessage.Index].decision = Value{
		id:  majorityValue.id,
		fit: majorityValue.fit,
	}
	in.proposerReplicatedLog[consensusMessage.Index].proposer = in.getProposer(majorityValue)

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
		in.debug("sending a decide consensus message to " + strconv.Itoa(int(i)))

		in.sendMessage(i, rpcPair)
	}
}

/*
	Proposer: Scan the E set and assign the number of times the Hi priority appears for each proposal
*/

func (in *Instance) getProposalWithMajorityHi(e []Value) Value {
	var ValueCount map[string]int64 // maps each id to hi count
	ValueCount = make(map[string]int64)
	for i := 0; i < len(e); i++ {
		if strings.HasPrefix(e[i].fit, strconv.FormatInt(in.Hi, 10)) {
			_, ok := ValueCount[e[i].id]
			if !ok {
				ValueCount[e[i].id] = 1 // note that each proposal has a unique id
			} else {
				ValueCount[e[i].id] = ValueCount[e[i].id] + 1
			}
		}
	}
	// check whether there is a value who has f+1 Hi count

	for key, element := range ValueCount {
		if element >= in.numReplicas/2+1 {
			for i := 0; i < len(e); i++ {
				if e[i].id == key {
					return e[i]
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
	// note that for each replica a different RPC pair should be generated

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
		in.debug("sending a generic consensus propose message to " + strconv.Itoa(int(i)))

		in.sendMessage(i, rpcPair)
	}
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
		in.debug("Re proposed because the decided value is not the one I proposed earlier ")
	}
}

/*
	invoked whenever a new slot is decided
	marks the slot as committed
	sends the response back to the client
    //todo how does pipeline affect this? when failures occur slot i will be completed before slot i-n, in that case should we implement a new failure detector and repurpose?
*/

func (in *Instance) updateStateMachine() {
	for in.proposerReplicatedLog[in.committedIndex+1].decided == true {
		decision := in.proposerReplicatedLog[in.committedIndex+1].decision.id
		messageBlock, ok := in.messageStore.Get(decision)
		if !ok {
			/*
				I haven't received the actual message block still, request it
			*/
			in.sendMessageBlockRequest(decision)
			return
		}

		in.proposerReplicatedLog[in.committedIndex+1].committed = true
		in.debug("Committed " + strconv.Itoa(int(in.committedIndex+1)))
		in.committedIndex++
		in.executeAndSendResponse(messageBlock)
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
			in.debug("Sent client response batch to " + strconv.Itoa(int(clientRequestBatch.Sender)))
		}

	}
}

/*
	Send a consensus request to the leader / set of leaders. Upon receiving this, the leader node will eventually propose a value for a slot, for this hash
*/

func (in *Instance) sendConsensusRequest(hash string) {
	leader := in.getDeterministicLeader1()
	consensusRequest := proto.ConsensusRequest{
		Sender:   in.nodeName,
		Receiver: leader,
		Hash:     hash,
	}
	rpcPair := RPCPair{
		Code: in.consensusRequestRpc,
		Obj:  &consensusRequest,
	}
	in.debug("sending a consensus request to " + strconv.Itoa(int(leader)))

	in.sendMessage(leader, rpcPair)
}

/*
	Upon receiving a consensus request, the leader node will propose it to the consensus layer
*/

func (in *Instance) handleConsensusRequest(request *proto.ConsensusRequest) {
	// todo check pipeline length
	if in.nodeName == in.getDeterministicLeader1() {
		in.updateProposedIndex(request.Hash)
		in.propose(in.proposedIndex, request.Hash)
		in.debug("Proposed " + request.Hash + " to " + strconv.Itoa(int(in.proposedIndex)))
	}
}

/*
	Set proposed index = max(proposed_index, committed index)+1
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
	in.debug("Added " + hash + " to proposed array position " + strconv.Itoa(int(in.proposedIndex)))
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
