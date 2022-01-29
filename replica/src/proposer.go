package raxos

import (
	"raxos/proto"
	"strconv"
)

/*
	This is a message from the recorder to the proposer
*/

func (in *Instance) handleProposerConsensusMessage(consensus *proto.GenericConsensus) {
	//todo implementation
}

/*
	propose a new block for the slot index
*/

func (in *Instance) propose(index int64, hash string) {
	// propose hash for the slot[index]
}

/*
	indication from the consensus layer that a value is decided
*/

func (in *Instance) delivered(index int, hash string, proposer int64) {
	in.updateStateMachine()
	if len(in.proposed) > index && in.proposed[index] != "" && in.proposed[index] != hash {
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
