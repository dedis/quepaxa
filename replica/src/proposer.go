package raxos

import "raxos/proto"

func (in *Instance) handleProposerConsensusMessage(consensus *proto.GenericConsensus) {

}

func (in *Instance) propose(index int, hash string) {
	// propose hash for the slot[index]
}

func (in *Instance) delivered(index int, hash string, proposer int64) {
	// indication from the consensus layer

}

func (in *Instance) updateStateMachine() {
	// while true loop
}
