package raxos

import (
	"log"
	"math/rand"
	"os"
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
				code: in.messageBlockRpc,
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
	messageBlockAck := proto.MessageBlockAck{
		Sender:   in.nodeName,
		Receiver: block.Sender,
		Hash:     block.Hash,
	}

	rpcPair := RPCPair{
		code: in.messageBlockAckRpc,
		Obj:  &messageBlockAck,
	}

	in.sendMessage(block.Sender, rpcPair)

}

func (in *Instance) handleMessageBlockRequest(request *proto.MessageBlockRequest) {
	messageBlock, ok := in.messageStore.Get(request.Hash)
	if ok {
		// the block exists
		messageBlock.Sender = in.nodeName
		messageBlock.Receiver = request.Sender
		rpcPair := RPCPair{
			code: in.messageBlockRpc,
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
		code: in.messageBlockRequestRpc,
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
	if request.Operation == 1 {
		in.startServer()
	} else if request.Operation == 2 {
		in.printLog()
	}
}

func (in *Instance) sendSampleClientResponse(ack *proto.MessageBlockAck) {
	messageBlock, _ := in.messageStore.Get(ack.Hash)
	// for each client block (client batch of requests) create a client batch reposne, and send it with the correct id

	for i := 0; i < len(messageBlock.Requests); i++ {
		clientRequestBatch := messageBlock.Requests[i]
		var replies []*proto.ClientResponseBatch_SingleClientResponse
		for j := 0; j < len(clientRequestBatch.Requests); j++ {
			replies = append(replies, &proto.ClientResponseBatch_SingleClientResponse{Message: clientRequestBatch.Requests[j].Message})
		}
		responseBatch := proto.ClientResponseBatch{
			Sender:    in.nodeName,
			Receiver:  clientRequestBatch.Sender,
			Responses: replies,
			Id:        clientRequestBatch.Id,
		}
		rpcPair := RPCPair{
			code: in.clientResponseBatchRpc,
			Obj:  &responseBatch,
		}

		in.sendMessage(clientRequestBatch.Sender, rpcPair)

	}
}

func (in *Instance) handleMessageBlockAck(ack *proto.MessageBlockAck) {
	in.messageStore.addAck(ack.Hash)
	acks := in.messageStore.getAcks(ack.Hash)
	if acks != nil && int64(len(acks)) == in.numReplicas/2+1 {
		// todo this block is guaranteed to be present in f+1 replicas, so its persistent
		// todo invoke consensus or send to leader
		// todo remove the following which are sample client responses
		in.sendSampleClientResponse(ack)
	}
}

func (in *Instance) handleClientStatusResponse(response *proto.ClientStatusResponse) {
	// replica doesn't receive a client status response
}

func (in *Instance) startServer() {

	go in.waitForConnections()
	time.Sleep(2 * time.Second)

	in.connectToReplicas()
	time.Sleep(2 * time.Second)

	in.startConnectionListners()
	time.Sleep(2 * time.Second)

	in.startOutgoingLinks()
	time.Sleep(2 * time.Second)

	in.run()
	time.Sleep(2 * time.Second)

	in.broadcastBlock()
	time.Sleep(2 * time.Second)

	in.updateStateMachine()
	time.Sleep(2 * time.Second)

}

func (in *Instance) printLog() {
	f, err := os.Create(in.logFilePath + strconv.Itoa(int(in.nodeName)) + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	choiceNum := 0
	for _, entry := range in.replicatedLog {
		// a single entry contains a 2D sequence of commands
		if entry.committed {
			// if an entry is committed, then it should contain the block in the message store
			block, ok := in.messageStore.Get(entry.decision.id)
			choiceLocalNum := 0
			if ok {
				for i := 0; i < len(block.Requests); i++ {
					for j := 0; j < len(block.Requests[i].Requests); j++ {
						_, _ = f.WriteString(strconv.Itoa(choiceNum) + "." + strconv.Itoa(choiceLocalNum) + ":")
						_, _ = f.WriteString(block.Requests[i].Requests[j].Message + ",")
						choiceLocalNum++
					}
				}
			} else {
				_, _ = f.WriteString(strconv.Itoa(choiceNum) + "." + strconv.Itoa(choiceLocalNum) + ":")
				_, _ = f.WriteString("no-op" + ",")
			}

		} else {
			_, _ = f.WriteString(strconv.Itoa(choiceNum) + "." + strconv.Itoa(0) + ":")
			_, _ = f.WriteString("no-op" + ",")
		}
		choiceNum++
	}
}
