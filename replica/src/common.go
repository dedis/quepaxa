package raxos

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"raxos/proto"
	"strconv"
	"strings"
	"time"
)

/*
	common.go implements the methods/functions that are common to both the proposer and the recorder
    	1. handling the hash overlay
		2. handling the consensus decisions, and log slot creation
*/

/*
	handler for new client requests received, the requests are sent to a channel for batching
	we allow clients to send requests in an open loop with an arbitrary passion arrival rate. To avoid chanel blocking, some client requests will be dropped
*/

func (in *Instance) handleClientRequestBatch(batch *proto.ClientRequestBatch) {

	// forward the batch of client requests to the buffer
	select {
	case in.requestsIn <- batch:
		// Success: the server side buffers are not full
		in.debug("successful pushing into server batching", 0)
	default:
		//Unsuccessful
		// if the buffer is full, then this request will be dropped (failed request)
		in.debug("unsuccessful pushing into server batching", 2)
	}

}

/*
	this is an infinite loop
	it collects a batch of client requests batches (a 2d array of requests), creates a new block and broadcasts it to all the replicas
*/

func (in *Instance) BroadcastBlock() {

	go func() {
		lastSent := time.Now() // used to get how long to wait
		for true {
			numRequests := int64(0)
			var requests []*proto.ClientRequestBatch
			// this loop collects requests until the minimum batch time is met OR the batch time is timeout
			for !(numRequests >= in.batchSize || (time.Now().Sub(lastSent).Microseconds() > in.batchTime && numRequests > 0)) {
				newRequest := <-in.requestsIn // keep collecting new requests for the next batch
				requests = append(requests, newRequest)
				numRequests++
			}

			in.blockCounter++

			// create a new mem block
			messageBlock := proto.MessageBlock{
				Sender:   in.nodeName,
				Receiver: 0,
				Hash:     strconv.Itoa(int(in.nodeName)) + "." + strconv.Itoa(int(in.blockCounter)), // unique block sequence number
				Requests: in.convertToMessageBlockRequests(requests),
			}

			// add the new message block to the store
			in.messageStore.Add(&messageBlock)

			// broadcast a new mem pool block
			for i := int64(0); i < in.numReplicas; i++ {
				messageBlock := proto.MessageBlock{
					Sender:   messageBlock.Sender,
					Receiver: i,
					Hash:     messageBlock.Hash, // unique block sequence number
					Requests: messageBlock.Requests,
				}
				rpcPair := RPCPair{
					Code: in.messageBlockRpc,
					Obj:  &messageBlock,
				}
				in.sendMessage(i, rpcPair)
			}
			in.debug("Sent a new mem block with hash "+messageBlock.Hash+" and with batch size "+strconv.Itoa(len(requests)), 0)
			lastSent = time.Now()
		}

	}()

}

/*
	when a new message block is received it is added to the store, and an ack is sent back to the creator of the block
*/

func (in *Instance) handleMessageBlock(block *proto.MessageBlock) {
	// add this block to the MessageStore
	in.messageStore.Add(block) // duplicates are handled internally
	in.debug("Added a message block from "+strconv.Itoa(int(block.Sender)), 0)
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
	in.debug("Send block ack to "+strconv.Itoa(int(block.Sender)), 0)
}

/*
	handler for message block acks. when a replica broadcasts a new MessageBlock, it retrieves acks.
	when a replica receives f+1 acks, that means that this block is persistent (less than f+1 replicas can fail by assumption)
	when it receives f+1 acks, the replica sends a consensus request message to the leader / sequence of leaders
*/

func (in *Instance) handleMessageBlockAck(ack *proto.MessageBlockAck) {
	in.messageStore.addAck(ack.Hash)
	acks := in.messageStore.getAcks(ack.Hash)
	if acks != nil && int64(len(acks)) == in.numReplicas/2+1 {
		in.debug("-----Received majority block acks for ---"+ack.Hash, 0)
		// note that this block is guaranteed to be present in f+1 replicas, so its persistent
		//in.sendSampleClientResponse(ack) //  this is only to test the overlay
		in.hashProposalsIn <- ack.Hash
	}
}

/*
	Upon receiving a message block request, send the requested block if it is found in the local message store
*/

func (in *Instance) handleMessageBlockRequest(request *proto.MessageBlockRequest) {
	messageBlock, ok := in.messageStore.Get(request.Hash)
	if ok {
		// the block exists
		in.debug("Sending the requested block to "+strconv.Itoa(int(request.Sender)), 0)
		messageBlock.Receiver = request.Sender // note that we do not alter the message.sender field, because its used to determine who originated the block
		rpcPair := RPCPair{
			Code: in.messageBlockRpc,
			Obj:  messageBlock,
		}
		in.sendMessage(request.Sender, rpcPair)
	}
}

/*
	invoked when the replica needs the message block to commit a request, and it is not already available in the store
*/

func (in *Instance) sendMessageBlockRequest(hash string) {
	// send a Message block request to a random recorder

	randomPeer := int64(rand.Intn(int(in.numReplicas)))
	for randomPeer == in.nodeName {
		randomPeer = int64(rand.Intn(int(in.numReplicas))) // to avoid sending to self
	}
	messageBlockRequest := proto.MessageBlockRequest{Hash: hash, Sender: in.nodeName, Receiver: int64(randomPeer)}
	rpcPair := RPCPair{
		Code: in.messageBlockRequestRpc,
		Obj:  &messageBlockRequest,
	}
	in.debug("sending a message block request to "+strconv.Itoa(int(randomPeer)), 0)

	in.sendMessage(int64(randomPeer), rpcPair)

}

/*
	append slots upto index 'index'
*/

func (in *Instance) initializeSlot(log []Slot, index int64) []Slot {

	for i := int64(len(log)); i < index+1000; i++ {
		log = append(log, Slot{
			index:                   i,
			proposedValue:           "",
			S:                       -1,
			P:                       []*proto.GenericConsensusValue{},
			E:                       []*proto.GenericConsensusValue{},
			C:                       []*proto.GenericConsensusValue{},
			committed:               false,
			decided:                 false,
			decision:                &proto.GenericConsensusValue{},
			proposer:                -1,
			proposeResponses:        []*proto.GenericConsensus{},
			spreadEResponses:        []*proto.GenericConsensus{},
			spreadCGatherEResponses: []*proto.GenericConsensus{},
			gatherCResponses:        []*proto.GenericConsensus{},
		})
	}
	return log
}

/*
	decision message is applicable to both proposer and the recorder
	upon receiving a decision message, both proposer log and the recorder log should be updated
*/

func (in *Instance) handleDecision(consensusMessage *proto.GenericConsensus) {

	in.recorderReplicatedLog = in.initializeSlot(in.recorderReplicatedLog, consensusMessage.Index) // create the slot if not already created
	in.proposerReplicatedLog = in.initializeSlot(in.proposerReplicatedLog, consensusMessage.Index) // create the slot if not already created

	if consensusMessage.D == true && in.proposerReplicatedLog[consensusMessage.Index].decided == false {
		in.recordProposerDecide(consensusMessage)
	}
	if consensusMessage.D == true && in.recorderReplicatedLog[consensusMessage.Index].decided == false {
		in.recordRecorderDecide(consensusMessage)
	}
}

/*
	// handler for generic consensus messages, corresponding method is called depending on the destination. Note that a replica acts as both a recorder and proposer
*/

func (in *Instance) handleGenericConsensus(consensus *proto.GenericConsensus) {

	if consensus.Destination == in.consensusMessageProposerDestination {
		in.handleProposerConsensusMessage(consensus)
	} else if consensus.Destination == in.consensusMessageRecorderDestination {
		in.handleRecorderConsensusMessage(consensus)
	} else if consensus.Destination == in.consensusMessageCommonDestination {
		in.handleDecision(consensus)
	}

}

/*
	Clients send status requests to (1) bootstrap and (2) print logs
*/

func (in *Instance) handleClientStatusRequest(request *proto.ClientStatusRequest) {
	if request.Operation == 1 {
		in.startServer()
	} else if request.Operation == 2 {
		in.printLog()
		//in.messageStore.printStore(in.logFilePath, in.nodeName) //this is for the testing purposes
	}

	/*send a status response back to the client*/
	statusResponse := proto.ClientStatusResponse{
		Sender:    in.nodeName,
		Receiver:  request.Sender,
		Operation: request.Operation,
		Message:   "Status Response from " + strconv.Itoa(int(in.nodeName)),
	}
	rpcPair := RPCPair{
		Code: in.clientStatusResponseRpc,
		Obj:  &statusResponse,
	}

	in.sendMessage(request.Sender, rpcPair)
	in.debug("Sent status response", 0)

}

/*
	this is a dummy method that is used for the testing of the message overlay. this method is triggered when a message block collects f+1 number of acks.
	In this implementation, for each client request batch, a response batch is generated that contains the same set of messages as the requests. Then for each batch of client requests
	a response is sent as a response batch.
*/

func (in *Instance) sendSampleClientResponse(ack *proto.MessageBlockAck) {
	messageBlock, _ := in.messageStore.Get(ack.Hash) // at this point the request is definitely in the store, so we don't check the existence
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
		in.debug("Sent client response batch to "+strconv.Itoa(int(clientRequestBatch.Sender)), 0)

	}
}

/*
	Server bootstrapping: 	Connect to other replicas
*/

func (in *Instance) startServer() {
	if !in.serverStarted {
		in.serverStarted = true
		in.connectToReplicas()
	}
}

/*
	print the replicated log, this is a util function that is used to check the log consistency
*/

func (in *Instance) printLog() {
	f, err := os.Create(in.logFilePath + strconv.Itoa(int(in.nodeName)) + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	choiceNum := 0
	for _, entry := range in.proposerReplicatedLog {
		// a single entry contains a 2D sequence of commands
		if entry.committed {
			// if an entry is committed, then it should contain the blocks in the message store
			decision := entry.decision.Id
			hashes := strings.Split(decision, ":")
			for i := 0; i < len(hashes); i++ {
				hashes[i] = strings.TrimSpace(hashes[i])
			}
			for i := 0; i < len(hashes); i++ {
				if len(hashes[i]) == 0 {
					continue
				}

				block, ok := in.messageStore.Get(hashes[i])
				if ok {
					for k := 0; k < len(block.Requests); k++ {
						for j := 0; j < len(block.Requests[k].Requests); j++ {
							_, _ = f.WriteString(strconv.Itoa(choiceNum) + "." + hashes[i] + "." + block.Requests[k].Id + ":")
							_, _ = f.WriteString(block.Requests[k].Requests[j].Message + "\n")
						}
					}
				} else {
					panic("committed entry " + fmt.Sprintf(" %v", entry) + " has mem blocks that are not found in the store")
				}
			}

		} else {
			break
		}
		choiceNum++
	}
}

/*
	this is an infinite thread
	it collects hashes of the blocks which are ready to be sent to the leader node
*/

func (in *Instance) CollectAndProposeHashes() {

	go func() {
		lastSent := time.Now() // used to get how long to wait
		for true {             // this runs forever
			numRequests := int(0)
			requestString := ""
			// this loop collects requests until the minimum batch time is met OR the batch time is timeout
			for !(numRequests > in.hashBatchSize || (time.Now().Sub(lastSent).Microseconds() > int64(in.hashBatchTime) && numRequests > 0)) {
				newRequest := <-in.hashProposalsIn // keep collecting new requests for the next batch
				requestString = requestString + ":" + newRequest
				numRequests++
			}

			in.debug("Collected a batch of hashes to propose "+requestString, 1)

			leader := in.getDeterministicLeader1() //todo change this when adding the optimizations

			consensusRequest := proto.ConsensusRequest{
				Sender:   in.nodeName,
				Receiver: leader,
				Hash:     requestString,
			}
			rpcPair := RPCPair{
				Code: in.consensusRequestRpc,
				Obj:  &consensusRequest,
			}
			in.debug("sending a consensus request to "+strconv.Itoa(int(leader))+" for the hash "+requestString, 0)

			in.sendMessage(leader, rpcPair)

			lastSent = time.Now()
		}
	}()
}