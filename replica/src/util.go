package raxos

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"raxos/configuration"
	"raxos/proto"
)

/*
	If enabled, print the messages to stdout
*/

func (in *Instance) debug(message string, level int) {
	if in.debugOn {
		if in.debugLevel <= level {
			fmt.Printf("%s\n", message)
		}
	}
}

/*
	Returns the size of any type of object in bytes
*/

func getRealSizeOf(v interface{}) (int, error) {
	b := new(bytes.Buffer)
	if err := gob.NewEncoder(b).Encode(v); err != nil {
		return 0, err
	}
	return b.Len(), nil
}

/*
	returns a deterministic string of size @length
*/

func getStringOfSizeN(length int) string {
	str := "a"
	size, _ := getRealSizeOf(str)
	for size < length {
		str = str + "a"
		size, _ = getRealSizeOf(str)
	}
	return str
}

func (in *Instance) convertClientRequest(request *proto.ClientRequestBatch_SingleClientRequest) *proto.MessageBlock_SingleClientRequest {
	var returnClientRequest proto.MessageBlock_SingleClientRequest
	returnClientRequest.Message = request.Message
	return &returnClientRequest
}

func (in *Instance) convertToClientRequestBatch(batch *proto.ClientRequestBatch) *proto.MessageBlock_ClientRequestBatch {
	var returnBatch proto.MessageBlock_ClientRequestBatch
	returnBatch.Id = batch.Id
	returnBatch.Sender = batch.Sender
	for i := 0; i < len(batch.Requests); i++ {
		returnBatch.Requests = append(returnBatch.Requests, in.convertClientRequest(batch.Requests[i]))
	}
	return &returnBatch
}

func (in *Instance) convertToMessageBlockRequests(requests []*proto.ClientRequestBatch) []*proto.MessageBlock_ClientRequestBatch {
	var returnArray []*proto.MessageBlock_ClientRequestBatch
	for i := 0; i < len(requests); i++ {
		clientRequestBatch := requests[i]
		returnArray = append(returnArray, in.convertToClientRequestBatch(clientRequestBatch))
	}
	return returnArray
}

/*
	Checks if this block was previous proposed
*/

func (in *Instance) proposedPreviously(hash string) (bool, int) {
	// checks if this value appears in the proposed array
	for i := 0; i < len(in.proposed); i++ {
		if in.proposed[i] == hash {
			return true, i
		}
	}
	return false, -1
}

/*
	checks of this value has been previously committed
*/

func (in *Instance) committedPreviously(hash string) (bool, int) {
	// checks if this value is previously decided
	for i := 0; i < len(in.proposerReplicatedLog); i++ {
		if in.proposerReplicatedLog[i].decision.id == hash {
			return true, i
		}
	}
	return false, -1
}

/*
	returns a fixed leader (strawman 1)
*/

func (in *Instance) getDeterministicLeader1() int64 {
	return 0 // node 0 is the default leader
}

/*
	extracts clients ip:port list to an array
*/

func getClientAddressList(cfg *configuration.InstanceConfig) []string {
	var clients []string
	for i := 0; i < len(cfg.Clients); i++ {
		clients = append(clients, cfg.Clients[i].Address)
	}
	return clients
}

/*
	extracts replicas ip:port list to an array
*/

func GetReplicaAddressList(cfg *configuration.InstanceConfig) []string {
	var replicas []string
	for i := 0; i < len(cfg.Peers); i++ {
		replicas = append(replicas, cfg.Peers[i].Address)
	}
	return replicas
}

/*
	Util: Convert data structure value array to proto generic consensus value array
*/

func (in *Instance) getGenericConsensusValueArray(e []Value) []*proto.GenericConsensusValue {
	var returnArray []*proto.GenericConsensusValue
	for i := 0; i < len(e); i++ {
		returnArray = append(returnArray, &proto.GenericConsensusValue{
			Id:  e[i].id,
			Fit: e[i].fit,
		})
	}
	return returnArray
}

/*
	Util: Convert generic consensus value array to []value
*/

func (in *Instance) proposerConvertToValueArray(c []*proto.GenericConsensusValue) []Value {
	returnArray := make([]Value, 0)
	for i := 0; i < len(c); i++ {
		returnArray = append(returnArray, Value{
			id:  c[i].Id,
			fit: c[i].Fit,
		})
	}
	return returnArray
}
