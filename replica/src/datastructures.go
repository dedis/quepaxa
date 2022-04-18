package raxos

import "raxos/proto"

/*
	Slot maintains the consensus and the smr state of a single Slot
	the definition of the sets is different from the algorithm, but the only different is the naming of sets. The core logic is same
*/
type Slot struct {
	index     int64
	S         int64                        // step
	P         *proto.GenericConsensusValue // proposed value
	E         []*proto.GenericConsensusValue
	C         []*proto.GenericConsensusValue
	U         []*proto.GenericConsensusValue
	committed bool
	decided   bool
	decision  *proto.GenericConsensusValue
	proposer  int64 // id of the proposer who decided this index

	proposeResponses        []*proto.GenericConsensus
	spreadEResponses        []*proto.GenericConsensus
	spreadCGatherEResponses []*proto.GenericConsensus
	gatherCResponses        []*proto.GenericConsensus
}
