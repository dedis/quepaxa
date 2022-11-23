package common

import (
	"bytes"
	"encoding/gob"
	"raxos/proto"
	"time"
)

/*
	RPC pair assigns a unique id to each type of message defined in the proto files
*/

type RPCPair struct {
	Code uint8
	Obj  proto.Serializable
}

/*
	Outgoing RPC assigns a rpc to its intended destination peer, the peer can be a replica or a client
*/

type OutgoingRPC struct {
	RpcPair *RPCPair
	Peer    int64
}

/*
	Timer for triggering event upon timeout
*/

type TimerWithCancel struct {
	d time.Duration
	t *time.Timer
	c chan interface{}
	f func()
}

/*
	instantiate a new timer with cancel
*/

func NewTimerWithCancel(d time.Duration) *TimerWithCancel {
	t := &TimerWithCancel{}
	t.d = d
	t.c = make(chan interface{}, 5)
	return t
}

/*
	Start the timer
*/

func (t *TimerWithCancel) Start() {
	t.t = time.NewTimer(t.d)
	go func() {
		select {
		case <-t.t.C:
			t.f()
			break
		case <-t.c:
			break
		}
	}()
}

/*
	Set a function to call when timeout
*/

func (t *TimerWithCancel) SetTimeoutFuntion(f func()) {
	t.f = f
}

/*
	Cancel timer
*/
func (t *TimerWithCancel) Cancel() {
	select {
	case t.c <- nil:
		// Success
		break
	default:
		//Unsuccessful
		break
	}

}

/*
	A util function to get the size of a message
*/

func GetRealSizeOf(v interface{}) (int, error) {
	b := new(bytes.Buffer)
	if err := gob.NewEncoder(b).Encode(v); err != nil {
		return 0, err
	}
	return b.Len(), nil
}