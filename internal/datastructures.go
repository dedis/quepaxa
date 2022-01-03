package internal

/*
	client_request corresponds to a single client request
*/
type client_request struct {
	id      string // assumed to be unique
	message string
}

/*
	message_block corresponds to a batch of requests that will be replicated
*/
type message_block struct {
	requests []client_request
	id       string // assumed to be unique
}

/*
	value associates the id of message block with a priority assigned by the recorders
*/

type value struct {
	id  string // id of the message_block
	fit int
}

/*
	slot maintains the consensus and the smr state of a single slot
*/
type slot struct {
	index     int
	S         int // step
	P         value
	E         []value
	C         []value
	U         []value
	committed bool
	decided   bool
	decision  value
	proposer  int // id of the proposer who decided this index
}
