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

func (in *Instance) debug(message string) {
	if in.debugOn {
		fmt.Printf("%s\n", message)
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
Creates a new copy of the message. Since the protobuf methods are not thread safe, its not possible to broadcast the same message without having separate message object for each
*/

func (in *Instance) getNewCopyOfMessage(code uint8, msg proto.Serializable) proto.Serializable {

	if code == in.clientRequestBatchRpc {

		clientRequestBatch := msg.(*proto.ClientRequestBatch)
		return &proto.ClientRequestBatch{
			Sender:   clientRequestBatch.Sender,
			Receiver: clientRequestBatch.Receiver,
			Requests: clientRequestBatch.Requests,
			Id:       clientRequestBatch.Id,
		}

	}

	if code == in.clientResponseBatchRpc {
		clientResponseBatch := msg.(*proto.ClientResponseBatch)
		return &proto.ClientResponseBatch{
			Sender:    clientResponseBatch.Sender,
			Receiver:  clientResponseBatch.Receiver,
			Responses: clientResponseBatch.Responses,
			Id:        clientResponseBatch.Id,
		}

	}

	if code == in.clientStatusRequestRpc {
		clientStatusRequest := msg.(*proto.ClientStatusRequest)
		return &proto.ClientStatusRequest{
			Sender:    clientStatusRequest.Sender,
			Receiver:  clientStatusRequest.Receiver,
			Operation: clientStatusRequest.Operation,
			Message:   clientStatusRequest.Message,
		}
	}

	if code == in.clientStatusResponseRpc {
		clientStatusResponse := msg.(*proto.ClientStatusResponse)
		return &proto.ClientStatusResponse{
			Sender:    clientStatusResponse.Sender,
			Receiver:  clientStatusResponse.Receiver,
			Operation: clientStatusResponse.Operation,
			Message:   clientStatusResponse.Message,
		}
	}

	if code == in.genericConsensusRpc {

		genericConsensus := msg.(*proto.GenericConsensus)
		return &proto.GenericConsensus{
			Sender:      genericConsensus.Sender,
			Receiver:    genericConsensus.Receiver,
			Index:       genericConsensus.Index,
			M:           genericConsensus.M,
			S:           genericConsensus.S,
			P:           genericConsensus.P,
			E:           genericConsensus.E,
			C:           genericConsensus.C,
			D:           genericConsensus.D,
			DS:          genericConsensus.DS,
			PR:          genericConsensus.PR,
			Destination: genericConsensus.Destination,
		}

	}

	if code == in.messageBlockRpc {

		messageBlock := msg.(*proto.MessageBlock)
		return &proto.MessageBlock{
			Sender:   messageBlock.Sender,
			Receiver: messageBlock.Receiver,
			Hash:     messageBlock.Hash,
			Requests: messageBlock.Requests,
		}

	}

	if code == in.messageBlockRequestRpc {
		messageBlockRequest := msg.(*proto.MessageBlockRequest)
		return &proto.MessageBlockRequest{
			Sender:   messageBlockRequest.Sender,
			Receiver: messageBlockRequest.Receiver,
			Hash:     messageBlockRequest.Hash,
		}
	}

	if code == in.messageBlockAckRpc {
		messageBlockAck := msg.(*proto.MessageBlockAck)
		return &proto.MessageBlockAck{
			Sender:   messageBlockAck.Sender,
			Receiver: messageBlockAck.Receiver,
			Hash:     messageBlockAck.Hash,
		}
	}

	return nil

}

func (in *Instance) convertClientRequest(request *proto.ClientRequestBatch_SingleClientRequest) *proto.MessageBlock_SingleClientRequest {
	var returnClientRequest *proto.MessageBlock_SingleClientRequest
	returnClientRequest.Message = request.Message
	return returnClientRequest
}

func (in *Instance) convertToClientRequestBatch(batch *proto.ClientRequestBatch) *proto.MessageBlock_ClientRequestBatch {
	var returnBatch *proto.MessageBlock_ClientRequestBatch
	returnBatch.Id = batch.Id
	returnBatch.Sender = batch.Sender
	for i := 0; i < len(batch.Requests); i++ {
		returnBatch.Requests = append(returnBatch.Requests, in.convertClientRequest(batch.Requests[i]))
	}
	return returnBatch
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
	for i := 0; i < len(in.replicatedLog); i++ {
		if in.replicatedLog[i].decision.id == hash {
			return true, i
		}
	}
	return false, -1
}

/*
 returns a fixed leader (strawman1)
*/

func (in *Instance) getDeterministicLeader1() int {
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

func getReplicaAddressList(cfg *configuration.InstanceConfig) []string {
	var replicas []string
	for i := 0; i < len(cfg.Peers); i++ {
		replicas = append(replicas, cfg.Peers[i].Address)
	}
	return replicas
}
