package common

import (
	"raxos/proto"
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
