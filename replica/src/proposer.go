package raxos

import (
	"context"
	"fmt"
	"math/rand"
	"raxos/proto/client"
	"sync"
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
	hi                       int
}

// instantiate a new Proposer

func NewProposer(name int64, threadId int64, peers []peer, proxyToProposerChan chan ProposeRequest, proposerToProxyChan chan ProposeResponse, proxyToProposerFetchChan chan FetchRequest, proposerToProxyFetchChan chan FetchResposne, lastSeenTimes []*time.Time, debugOn bool, debugLevel int, hi int) *Proposer {

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
		hi:                       hi,
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

// convert between proto types

func (prop *Proposer) convertToClientBtchMessages(messages []*client.ClientBatch_SingleMessage) []*ProposerMessage_ClientBatch_SingleMessage {
	rA := make([]*ProposerMessage_ClientBatch_SingleMessage, 0)
	for i := 0; i < len(messages); i++ {
		rA = append(rA, &ProposerMessage_ClientBatch_SingleMessage{
			Message: messages[i].Message,
		})
	}
}

// convert between proto types

func (prop *Proposer) getProposeClientBatches(btch []client.ClientBatch) []*ProposerMessage_ClientBatch {
	rA := make([]*ProposerMessage_ClientBatch, 0)
	for i := 0; i < len(btch); i++ {
		rA = append(rA, &ProposerMessage_ClientBatch{
			Sender:   btch[i].Sender,
			Messages: prop.convertToClientBtchMessages(btch[i].Messages),
			Id:       btch[i].Id,
		})
	}
}

// convert between proto types

func (prop *Proposer) extractDecidedSlots(indexes []int, decisions [][]string) []*ProposerMessage_DecidedSlot {
	if len(indexes) != len(decisions) {
		panic("should not happen")
	}
	rA := make([]*ProposerMessage_DecidedSlot, 0)
	for i := 0; i < len(indexes); i++ {
		rA = append(rA, &ProposerMessage_DecidedSlot{
			Index:    int64(indexes[i]),
			Ids:      decisions[i],
			Proposer: prop.name,
		})
	}

	return rA
}

// return the maximum of F’ from all replies in R

func (prop *Proposer) getMaxFromResponses(array []RecorderResponse, set string) ProposerMessage_Proposal {
	if len(array) == 0 {
		panic("should not happen")
	}
	
	var max *RecorderResponse_Proposal
	
	if set == "F" {
		max = array[0].F
	}else if set == "M"{
		max = array[0].M	
	}
	
	for i := 1; i < len(array); i++ {
		if prop.isGreaterThan(array[i], max, set) { //todo start from here
			if set == "F" {
				max = array[i].F
			}else if set == "M"{
				max = array[i].M
			}		
		}
	}
}

// run the proposer logic

func (prop *Proposer) handleProposeRequest(message ProposeRequest) ProposeResponse {
	S := 1*4 + 0
	P := ProposerMessage_Proposal{
		Priority:      int64(prop.hi),
		ProposerId:    prop.name,
		ThreadId:      prop.threadId,
		Ids:           message.proposalStr,
		ClientBatches: prop.getProposeClientBatches(message.proposalBtch),
	}

	decidedSlots := prop.extractDecidedSlots(message.lastDecidedIndexes, message.lastDecidedDecisions)

	for true {

		Pi := make([]ProposerMessage_Proposal, prop.numReplicas)

		//todo add fast path where the first try does not contain the client batches

		for i := 0; i < prop.numReplicas; i++ {
			Pi[i] = ProposerMessage_Proposal{
				Priority:      P.Priority,
				ProposerId:    P.ProposerId,
				ThreadId:      P.ThreadId,
				Ids:           P.Ids,
				ClientBatches: P.ClientBatches,
			}
		}
		if S%4 == 0 && (S > 4 || message.msWait != 0) {
			for i := 0; i < prop.numReplicas; i++ {
				Pi[i].Priority = int64(rand.Intn(prop.hi-1)) + 1
			}
		}

		responses := make(chan *RecorderResponse, prop.numReplicas)
		ctx, cancel := context.WithTimeout(context.Background(), in.timeout)

		wg := sync.WaitGroup{}
		for i := 0; i < prop.numReplicas; i++ {
			wg.Add(1)
			go func(p peer, pi ProposerMessage_Proposal, s int, decidedSlots []*ProposerMessage_DecidedSlot) {
				defer wg.Done()
				resp, err := p.client.ESP(ctx, &ProposerMessage{
					Sender:       prop.name,
					Index:        message.instance,
					P:            &pi,
					S:            int64(s),
					DecidedSlots: decidedSlots,
				})

				if err != nil {
					return
				}

				select {
				case responses <- resp:
					return
				default:
					return
				}

			}(prop.peers[i], Pi[i], S, decidedSlots)
		}
		decidedSlots = make([]*ProposerMessage_DecidedSlot, 0)

		go func() {
			wg.Wait()
			cancel()
		}()

		// todo add fast path checks where HasClientBacthes is false
		responsesArray := make([]RecorderResponse, 0)
		for r := range responses {
			responsesArray = append(responsesArray, *r)
			if len(responsesArray) == prop.numReplicas/2+1 {
				close(responses)
			}
		}
		// If all replies in R have S’ = S and S%4 = 0
		allRepliesHaveS := true
		for i := 0; i < len(responsesArray); i++ {
			if responsesArray[i].S != int64(S) {
				allRepliesHaveS = false
				break
			}
		}

		if allRepliesHaveS && S%4 == 0 {
			allRepliesHaveFHiFit := true
			for i := 0; i < len(responsesArray); i++ {
				if responsesArray[i].F.Priority != int64(prop.hi) {
					allRepliesHaveFHiFit = false
					break
				}
			}
			if allRepliesHaveFHiFit {
				return ProposeResponse{
					index:     int(message.instance),
					decisions: responsesArray[0].F.Ids,
				}
			}

			// P ← maximum of F’ from all replies in R
			P = prop.getMaxFromResponses(responsesArray, "F")
		}

	}
}

// infinite loop listening to the server channel

func (prop *Proposer) runProposer() {
	go func() {
		for true {

			select {
			case proposeMessage := <-prop.proxyToProposerChan:
				prop.debug("Received propose request", 0)
				prop.proposerToProxyChan <- prop.handleProposeRequest(proposeMessage)
				break

			case fetchMessage := <-prop.proxyToProposerFetchChan:
				prop.debug("Received fetch request", 0)
				prop.proposerToProxyFetchChan <- prop.handleFetchRequest(fetchMessage)
				break
			}
		}

	}()
}
