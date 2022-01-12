package raxos

import (
	"math/rand"
	"raxos/proto"
	"strconv"
	"time"
)

func (in *Instance) broadcastBlock() {
	go func() {
		lastSent := time.Now() // used to get how long to wait
		for true {             // this runs forever
			numRequests := int64(0)
			var requests []*proto.ClientRequestBatch
			for !(numRequests >= in.batchSize || (time.Now().Sub(lastSent).Microseconds() > in.batchTime && numRequests > 0)) {
				newRequest := <-in.requestsIn // keep collecting new requests for the next batch
				requests = append(requests, newRequest)
				numRequests++
			}

			messageBlock := proto.MessageBlock{
				Sender:   in.nodeName,
				Receiver: 0,
				Hash:     strconv.Itoa(int(in.nodeName)) + "." + strconv.Itoa(int(in.blockCounter)),
				Requests: in.convertToMessageBlockRequests(requests),
			}

			rpcPair := RPCPair{
				code: in.MessageBlockRpc,
				Obj:  &messageBlock,
			}

			for i := int64(0); i < in.numReplicas; i++ {
				in.sendMessage(i, rpcPair)
			}

			lastSent = time.Now()
		}

	}()

}

func (in *Instance) handleClientRequestBatch(batch *proto.ClientRequestBatch) {

	// forward the batch of client requests to the requests in buffer
	select {
	case in.requestsIn <- batch:
		// Success
	default:
		//Unsuccessful
		// if the buffer is full, then this request will be dropped (failed request)
	}

}

func (in *Instance) handleClientResponseBatch(batch *proto.ClientResponseBatch) {
	// the proposer doesn't receive any client responses
}

func (in *Instance) handleMessageBlock(block *proto.MessageBlock) {
	// add this block to the MessageStore
	in.messageStore.Add(block)

}

func (in *Instance) handleMessageBlockRequest(request *proto.MessageBlockRequest) {
	messageBlock, ok := in.messageStore.Get(request.Hash)
	if ok {
		// the block exists
		messageBlock.Sender = in.nodeName
		messageBlock.Receiver = request.Sender

		rpcPair := RPCPair{
			code: in.MessageBlockRpc,
			Obj:  messageBlock,
		}

		in.sendMessage(request.Sender, rpcPair)

	}
}

func (in *Instance) sendMessageBlockRequest(hash string) {
	// send a Message block request to a random recorder

	randomPeer := rand.Intn(int(in.numReplicas))
	messageBlockRequest := proto.MessageBlockRequest{Hash: hash, Sender: in.nodeName, Receiver: int64(randomPeer)}
	rpcPair := RPCPair{
		code: in.MessageBlockRequestRpc,
		Obj:  &messageBlockRequest,
	}

	in.sendMessage(int64(randomPeer), rpcPair)

}

func (in *Instance) handleGenericConsensus(consensus *proto.GenericConsensus) {
	// 1 for the proposer and 2 for the recorder
	if consensus.Destination == 1 {
		in.handleProposerConsensusMessage(consensus)
	} else if consensus.Destination == 2 {
		in.handleRecorderConsensusMessage(consensus)
	}

}

func (in *Instance) handleClientStatusRequest(request *proto.ClientStatusRequest) {

}

func (in *Instance) handleClientStatusResponse(response *proto.ClientStatusResponse) {

}
