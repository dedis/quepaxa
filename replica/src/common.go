package raxos

import (
	"log"
	"math/rand"
	"os"
	"raxos/proto"
	"strconv"
	"time"
)

/*
	common.go implements the methods/functions that are common to both the proposer and the recorder
*/

func (in *Instance) BroadcastBlock() {
	/*
		This is an infinite thread
		It collects a batch of client requests batches (a 2d array of requests), creates a new block and broadcasts it to all the replicas
	*/
	go func() {
		lastSent := time.Now() // used to get how long to wait
		for true {             // this runs forever
			numRequests := int64(0)
			var requests []*proto.ClientRequestBatch
			// this loop collects requests until the minimum batch time is met OR the batch time is timeout
			for !(numRequests >= in.batchSize || (time.Now().Sub(lastSent).Microseconds() > in.batchTime && numRequests > 0)) {
				newRequest := <-in.requestsIn // keep collecting new requests for the next batch
				requests = append(requests, newRequest)
				numRequests++
			}

			messageBlock := proto.MessageBlock{
				Sender:   in.nodeName,
				Receiver: 0,                                                                         // message block is a broadcast, so the receiver field is of no use
				Hash:     strconv.Itoa(int(in.nodeName)) + "." + strconv.Itoa(int(in.blockCounter)), // unique block sequence number
				Requests: in.convertToMessageBlockRequests(requests),
			}

			rpcPair := RPCPair{
				Code: in.messageBlockRpc,
				Obj:  &messageBlock,
			}

			for i := int64(0); i < in.numReplicas; i++ {
				in.sendMessage(i, rpcPair)
			}

			lastSent = time.Now()
		}

	}()

}

/*
handler for new client requests received, the requests are sent to a channel for batching
we allow clients to send requests in an open loop with an arbitrary passion arrival rate. To avoid buffer overflows, some client requests will be dropped
*/

func (in *Instance) handleClientRequestBatch(batch *proto.ClientRequestBatch) {

	// forward the batch of client requests to the requests in buffer
	select {
	case in.requestsIn <- batch:
		// Success: the server side buffers are not full
	default:
		//Unsuccessful
		// if the buffer is full, then this request will be dropped (failed request)
	}

}

/*
At the moment a replica does not receive client response batches
*/

func (in *Instance) handleClientResponseBatch(batch *proto.ClientResponseBatch) {
	// the proposer doesn't receive any client responses
}

/*
	When a new message block is received it is added to the store, and an ack is sent back to the creator of the block
*/

func (in *Instance) handleMessageBlock(block *proto.MessageBlock) {
	// add this block to the MessageStore
	in.messageStore.Add(block)
	messageBlockAck := proto.MessageBlockAck{
		Sender:   in.nodeName,
		Receiver: block.Sender,
		Hash:     block.Hash,
	}

	rpcPair := RPCPair{
		Code: in.messageBlockAckRpc,
		Obj:  &messageBlockAck,
	}

	in.sendMessage(block.Sender, rpcPair)
}

/*
	Upon receiving a message block request, send the requested block if its found in the local message store
*/

func (in *Instance) handleMessageBlockRequest(request *proto.MessageBlockRequest) {
	messageBlock, ok := in.messageStore.Get(request.Hash)
	if ok {
		// the block exists
		messageBlock.Receiver = request.Sender
		rpcPair := RPCPair{
			Code: in.messageBlockRpc,
			Obj:  messageBlock,
		}
		in.sendMessage(request.Sender, rpcPair)
	}
}

/*invoked when the replica needs the message block to commit a request, and its not already available in the store*/

func (in *Instance) sendMessageBlockRequest(hash string) {
	// send a Message block request to a random recorder

	randomPeer := rand.Intn(int(in.numReplicas))
	messageBlockRequest := proto.MessageBlockRequest{Hash: hash, Sender: in.nodeName, Receiver: int64(randomPeer)}
	rpcPair := RPCPair{
		Code: in.messageBlockRequestRpc,
		Obj:  &messageBlockRequest,
	}

	in.sendMessage(int64(randomPeer), rpcPair)

}

/*handler for generic consensus messages, corresponding method is called depending on the destination. Note that a replica acts as both a recorder and proposer*/

func (in *Instance) handleGenericConsensus(consensus *proto.GenericConsensus) {
	// 1 for the proposer and 2 for the recorder
	if consensus.Destination == 1 {
		in.handleProposerConsensusMessage(consensus)
	} else if consensus.Destination == 2 {
		in.handleRecorderConsensusMessage(consensus)
	}

}

/*Clients send status requests to (1) bootstrap and print logs and (2) check status of the replica*/

func (in *Instance) handleClientStatusRequest(request *proto.ClientStatusRequest) {
	if request.Operation == 1 {
		in.startServer()
	} else if request.Operation == 2 {
		in.printLog()
		// todo remove the message store printing, its only for testing purposes
		in.messageStore.printStore(in.logFilePath, in.nodeName)
	}
}

/*
this is a dummy method that is used for the testing of the message overlay. This method is triggered when a message block collects f+1 number of acks.
In this implementation, for each client request batch, a response batch is generated that contains the same set of messages as the requests. Then for each batch of client requests
a response is sent as a response batch.

In the actual Raxos implementation, once a replica collects f+1  acks for its block, that replica sends a consensus request message to a leader / sequence of leaders
Then that leader runs the consensus algorithm. When the leader decides on the message he does not have to send back the response to the client. Since we broadcast the
decide messages, eventually the replica who originated that block will update the state machine, and find the client who sent the block in the MessageBlock.
Then that client will be sent a reply

This approach decreases the overhead imposed on the leader node.

The drawback of this approach is, if the block originator (the replica which created the block is crashed, then the client will not receive the response, even though
it is committed. We consider these requests as failed, and they will appear as the error rate in the client metrics)
*/

func (in *Instance) sendSampleClientResponse(ack *proto.MessageBlockAck) {
	messageBlock, _ := in.messageStore.Get(ack.Hash) // at this point the request is defintly in the store, so we don't check the existence
	// for each client block (client batch of requests) create a client batch response, and send it with the correct id (unique id for a batch of client requests)

	for i := 0; i < len(messageBlock.Requests); i++ {
		clientRequestBatch := messageBlock.Requests[i]
		var replies []*proto.ClientResponseBatch_SingleClientResponse
		for j := 0; j < len(clientRequestBatch.Requests); j++ {
			replies = append(replies, &proto.ClientResponseBatch_SingleClientResponse{Message: clientRequestBatch.Requests[j].Message})
		}

		// for each client batch send a response batch

		responseBatch := proto.ClientResponseBatch{
			Sender:    in.nodeName,
			Receiver:  clientRequestBatch.Sender,
			Responses: replies,
			Id:        clientRequestBatch.Id,
		}
		rpcPair := RPCPair{
			Code: in.clientResponseBatchRpc,
			Obj:  &responseBatch,
		}

		in.sendMessage(clientRequestBatch.Sender, rpcPair)

	}
}

/*
Handler for message block acks. When a replica broadcasts a new MessageBlock, it retrieves acks.
When a replica receives f+1 acks, that means that this block is persistent (less than f+1 replicas can fail by assumption)
When it receives f+1 acks, the replica sends a consensus request message to the leader / sequence of leaders
*/

func (in *Instance) handleMessageBlockAck(ack *proto.MessageBlockAck) {
	in.messageStore.addAck(ack.Hash)
	acks := in.messageStore.getAcks(ack.Hash)
	if acks != nil && int64(len(acks)) == in.numReplicas/2+1 {
		// note that this block is guaranteed to be present in f+1 replicas, so its persistent
		// todo remove the following invocation is only for testing purposes of the overlay
		in.sendSampleClientResponse(ack)
		// todo invoke to send to leader
		in.sendConsensusRequest(ack.Hash)
	}
}

/*
Handler for client status responses, currently the replica does not receive client responses, only for testing purposes

*/

func (in *Instance) handleClientStatusResponse(response *proto.ClientStatusResponse) {
	// replica doesn't receive a client status response
}

/*
Server bootstrapping

1. Listen to incoming connections from replicas and clients
2. Connect to other replicas
3. Start connection listeners (tcp connections from clients and replicas)
4. Start outgoing links: a set of threads which writes to the wire
5. Start main thread that listens to all the incoming messages from the centralied channel
6. Start the client request batcher which collects and broadcasts new blocks

*/

func (in *Instance) startServer() {

	in.connectToReplicas()
	time.Sleep(2 * time.Second)

	in.startConnectionListners()
	time.Sleep(2 * time.Second)

}

/*
Print the replicated log, this is a util function that is used to check the log consistency
*/

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

			// in theory this execution path should not execute

			_, _ = f.WriteString(strconv.Itoa(choiceNum) + "." + strconv.Itoa(0) + ":")
			_, _ = f.WriteString("no-op" + ",")
		}
		choiceNum++
	}
}
