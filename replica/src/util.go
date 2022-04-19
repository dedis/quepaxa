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

/*
	a util function to convert a single client request between mem pool and client request
*/

func (in *Instance) convertClientRequest(request *proto.ClientRequestBatch_SingleClientRequest) *proto.MessageBlock_SingleClientRequest {
	var returnClientRequest proto.MessageBlock_SingleClientRequest
	returnClientRequest.Message = request.Message
	return &returnClientRequest
}

/*
	a util function to convert between client request and mem pool client request batch
*/

func (in *Instance) convertToClientRequestBatch(batch *proto.ClientRequestBatch) *proto.MessageBlock_ClientRequestBatch {
	var returnBatch proto.MessageBlock_ClientRequestBatch
	returnBatch.Id = batch.Id
	returnBatch.Sender = batch.Sender
	for i := 0; i < len(batch.Requests); i++ {
		returnBatch.Requests = append(returnBatch.Requests, in.convertClientRequest(batch.Requests[i]))
	}
	return &returnBatch
}

/*
	A util function to convert different client request batch array
*/

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
		if in.proposerReplicatedLog[i].decision.Id == hash {
			return true, i
		}
	}
	return false, -1
}

/*
	Util: returns a fixed leader (strawman 1)
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
	A util function to support set union, item wise, for proto value
*/

func (in *Instance) setUnionProtoValue(array []*proto.GenericConsensusValue, protoValue *proto.GenericConsensusValue) []*proto.GenericConsensusValue {
	for i := 0; i < len(array); i++ {
		if array[i].Id == protoValue.Id && array[i].Fit == protoValue.Fit {
			return array
		}
	}
	array = append(array, protoValue)
	return array
}

/*
	A util function to support set union, set wise, for proto value
*/

func (in *Instance) setUnionProtoValues(array []*proto.GenericConsensusValue, protoValue []*proto.GenericConsensusValue) []*proto.GenericConsensusValue {

	for i := 0; i < len(protoValue); i++ {
		found := false
		for j := 0; j < len(array); j++ {
			if array[j].Id == protoValue[i].Id && array[j].Fit == protoValue[i].Fit {
				found = true
				break
			}
		}
		if found == false {
			array = append(array, protoValue[i])
		}
	}

	return array
}

/*
	A util function to remove empty values from a slice
*/

func (in *Instance) removeEmptyValues(array []*proto.GenericConsensusValue) []*proto.GenericConsensusValue {
	returnArray := make([]*proto.GenericConsensusValue, 0)
	for i := 0; i < len(array); i++ {
		if array[i].Id != "" {
			returnArray = append(returnArray, array[i])
		} else {
			panic("Empty slot found")
		}
	}

	return returnArray
}
