package raxos

import "raxos/proto"

/*
	Slot maintains the consensus and the smr state of a single Slot
*/
type Slot struct {
	index         int64
	proposedValue string
	S             int64 // step
	P             []*proto.GenericConsensusValue
	E             []*proto.GenericConsensusValue
	C             []*proto.GenericConsensusValue
	committed     bool
	decided       bool
	decision      *proto.GenericConsensusValue
	proposer      int64 // id of the proposer who decided this index

	proposeResponses        []*proto.GenericConsensus
	spreadEResponses        []*proto.GenericConsensus
	spreadCGatherEResponses []*proto.GenericConsensus
	gatherCResponses        []*proto.GenericConsensus
}
