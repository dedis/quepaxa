package raxos

import "raxos/proto"

/*
	Value associates the id of message block with a priority assigned by the recorders
*/

type Value struct {
	id  string // id of the MessageBlock
	fit string
}

/*
	Slot maintains the consensus and the smr state of a single Slot
	the definition of the sets is different from the implementation, but the only different is the naming of sets. The core logic is same
*/
type Slot struct {
	index     int64
	S         int64 // step
	P         Value // proposed value
	E         []Value
	C         []Value
	U         []Value
	committed bool
	decided   bool
	decision  Value
	proposer  int64 // id of the proposer who decided this index

	proposeResponses        []*proto.GenericConsensus
	spreadEResponses        []*proto.GenericConsensus
	spreadCGatherEResponses []*proto.GenericConsensus
	gatherCResponses        []*proto.GenericConsensus
}
