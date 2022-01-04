package internal

/*
	ClientRequest corresponds to a single client request
*/
type ClientRequest struct {
	id      string // assumed to be unique
	message string
}

/*
	MessageBlock corresponds to a batch of requests that will be replicated
*/
type MessageBlock struct {
	id       string // assumed to be unique
	requests []ClientRequest
}

/*
	Value associates the id of message block with a priority assigned by the recorders
*/

type Value struct {
	id  string // id of the MessageBlock
	fit int
}

/*
	Slot maintains the consensus and the smr state of a single Slot
*/
type Slot struct {
	index     int
	S         int // step
	P         Value
	E         []Value
	C         []Value
	U         []Value
	committed bool
	decided   bool
	decision  Value
	proposer  int // id of the proposer who decided this index
}
