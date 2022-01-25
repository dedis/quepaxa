package raxos

import "raxos/proto"

/*
	This is a message from the recorder to the proposer
*/

func (in *Instance) handleProposerConsensusMessage(consensus *proto.GenericConsensus) {
	//todo implementation
}

/*
	propose a new block for the slot index
*/

func (in *Instance) propose(index int, hash string) {
	// propose hash for the slot[index]
}

/*
	indication from the consensus layer that a value is decided
*/

func (in *Instance) delivered(index int, hash string, proposer int64) {
	// indication from the consensus layer
}

/*
	invoked whenever a new slot is decided
	marks the slot as committed
	sends the response back to the client
    //todo how does pipeline affect this? when failures occur slot i will be completed before slot i-n, in that case should we implement a new failure detector and repurpose?
*/

func (in *Instance) updateStateMachine() {
	//

}

/*
	Send a consensus request to the leader / set of leaders. Upon receiving this, the leader node will eventually propose a value for a slot, for this hash
*/

func (in *Instance) sendConsensusRequest(hash string) {
	//todo implement
	// todo new consensus messages should be defined
}
