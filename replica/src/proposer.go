package raxos

import (
	"context"
	"fmt"
	"math/rand"
	"raxos/proto/client"
	"time"
)

type Proposer struct {
	numReplicas              int
	name                     int64
	threadId                 int64
	peers                    []peer // gRPC connection list
	proxyToProposerChan      chan ProposeRequest
	proposerToProxyChan      chan ProposeResponse
	proxyToProposerFetchChan chan FetchRequest
	proposerToProxyFetchChan chan FetchResposne
	lastSeenTimes            []*time.Time
	debugOn                  bool // if turned on, the debug messages will be print on the console
	debugLevel               int  // debug level
}

// instantiate a new Proposer

func NewProposer(name int64, threadId int64, peers []peer, proxyToProposerChan chan ProposeRequest, proposerToProxyChan chan ProposeResponse, proxyToProposerFetchChan chan FetchRequest, proposerToProxyFetchChan chan FetchResposne, lastSeenTimes []*time.Time, debugOn bool, debugLevel int) *Proposer {

	pr := Proposer{
		numReplicas:              len(peers),
		name:                     name,
		threadId:                 threadId,
		peers:                    peers,
		proxyToProposerChan:      proxyToProposerChan,
		proposerToProxyChan:      proposerToProxyChan,
		proxyToProposerFetchChan: proxyToProposerFetchChan,
		proposerToProxyFetchChan: proposerToProxyFetchChan,
		lastSeenTimes:            lastSeenTimes,
		debugOn:                  debugOn,
		debugLevel:               debugLevel,
	}

	return &pr
}

/*
	if turned on, print the message to console
*/

func (prop *Proposer) debug(message string, level int) {
	if prop.debugOn && level >= prop.debugLevel {
		fmt.Printf("%s\n", message)
	}
}

// return the grpc client of rn
func (prop *Proposer) findClientByName(rn int64) peer {
	for i := 0; i < len(prop.peers); i++ {
		if prop.peers[i].name == rn {
			return prop.peers[i]
		}
	}

	panic("replica not found")
}

// select a random grpc client that is not self
func (prop *Proposer) getRandomClient() peer {
	rn := rand.Intn(prop.numReplicas)
	for int64(rn) == prop.name {
		rn = rand.Intn(prop.numReplicas)
	}
	return prop.findClientByName(int64(rn))
}

// convert between proto types

func (prop *Proposer) convertToClientBatchMessages(messages []*DecideResponse_ClientBatch_SingleMessage) []*client.ClientBatch_SingleMessage {
	rtMessages := make([]*client.ClientBatch_SingleMessage, 0)
	for i := 0; i < len(messages); i++ {
		rtMessages = append(rtMessages, &client.ClientBatch_SingleMessage{
			Message: messages[i].Message,
		})
	}

	return rtMessages
}

// convert to client batches

func (prop *Proposer) convertToClientBatches(batches []*DecideResponse_ClientBatch) []client.ClientBatch {
	rtBtches := make([]client.ClientBatch, 0)
	for i := 0; i < len(batches); i++ {
		rtBtches = append(rtBtches, client.ClientBatch{
			Sender:   batches[i].Sender,
			Messages: prop.convertToClientBatchMessages(batches[i].Messages),
			Id:       batches[i].Id,
		})
	}
	return rtBtches
}

// send a fetch request to a random peer

func (prop *Proposer) handleFetchRequest(message FetchRequest) FetchResposne {

	found := false
	var cltBatches []*DecideResponse_ClientBatch
	for !found {

		client := prop.getRandomClient()
		clientCon := client.client
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(2*time.Second))
		resp, err := clientCon.FetchBatches(ctx, &DecideRequest{
			Ids: message.ids,
		})

		if err == nil && resp != nil {
			if len(resp.ClientBatches) == len(message.ids) {
				found = true
				cltBatches = resp.ClientBatches
			}

		}
		cancel()
	}

	response := FetchResposne{
		batches: prop.convertToClientBatches(cltBatches),
	}
	return response
}

// infinite loop listening to the server channel

func (prop *Proposer) runProposer() {
	go func() {
		for true {

			select {
			//case proposeMessage := <-prop.proxyToProposerChan:
			//	prop.debug("Received propose request", 0)
			//	prop.proposerToProxyChan <- prop.handleProposeRequest(proposeMessage)
			//	break

			case fetchMessage := <-prop.proxyToProposerFetchChan:
				prop.debug("Received fetch request", 0)
				prop.proposerToProxyFetchChan <- prop.handleFetchRequest(fetchMessage)
				break
			}
		}
		
	}()
}
