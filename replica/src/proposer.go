package raxos

import "raxos/proto"

func (in *Instance) handleClientRequest(request *proto.ClientRequest) {
	// forward the request to the buffer that collects batches of requests

}

func (in *Instance) handleClientResponse(response *proto.ClientResponse) {

}

func (in *Instance) handleMessageBlockReply(reply *proto.MessageBlockReply) {

}

func (in *Instance) handleMessageBlockRequest(request *proto.MessageBlockRequest) {

}
