package raxos

import (
	"context"
	"math/rand"
	"raxos/proto/client"
	"sync"
	"time"
)

// this file contains the proposer side logic of QuePaxa

type Proposer struct {
	numReplicas                 int
	name                        int64  // unique
	threadId                    int64  // unique
	peers                       []peer // gRPC connection list
	proxyToProposerChan         chan ProposeRequest
	proposerToProxyChan         chan ProposeResponse
	proxyToProposerFetchChan    chan FetchRequest
	proposerToProxyFetchChan    chan FetchResposne
	debugOn                     bool // if turned on, the debug messages will be print on the console
	debugLevel                  int  // debug level
	hi                          int  // hi priority
	serverMode                  int  // if 1, use the fast path optimizations
	proxyToProposerDecisionChan chan Decision
}

// instantiate a new Proposer

func NewProposer(name int64, threadId int64, peers []peer, proxyToProposerChan chan ProposeRequest, proposerToProxyChan chan ProposeResponse, proxyToProposerFetchChan chan FetchRequest, proposerToProxyFetchChan chan FetchResposne, debugOn bool, debugLevel int, hi int, serverMode int, proxyToProposerDecisionChan chan Decision) *Proposer {

	pr := Proposer{
		numReplicas:                 len(peers),
		name:                        name,
		threadId:                    threadId + 1,
		peers:                       peers,
		proxyToProposerChan:         proxyToProposerChan,
		proposerToProxyChan:         proposerToProxyChan,
		proxyToProposerFetchChan:    proxyToProposerFetchChan,
		proposerToProxyFetchChan:    proposerToProxyFetchChan,
		debugOn:                     debugOn,
		debugLevel:                  debugLevel,
		hi:                          hi,
		serverMode:                  serverMode,
		proxyToProposerDecisionChan: proxyToProposerDecisionChan,
	}

	//pr.debug("created a new proposer instance", -1)

	return &pr
}

///*
//	if turned on, print the message to console,
//*/
//
//func (prop *Proposer) debug(message string, level int) {
//	if prop.debugOn && level >= prop.debugLevel {
//		fmt.Printf("%s\n", message)
//	}
//}

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
	rn := rand.Intn(prop.numReplicas) + 1
	for int64(rn) == prop.name {
		rn = rand.Intn(prop.numReplicas) + 1
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
	//prop.debug("proposer starting to handle a fetch request "+message.ids[0], 1)
	found := false
	cltBatches := make([]*DecideResponse_ClientBatch, 0)
	numBtches := len(message.ids)

	for !found {

		client_r := prop.getRandomClient()
		clientCon := client_r.client
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
		resp, err := clientCon.FetchBatches(ctx, &DecideRequest{
			Ids: message.ids,
		})

		//prop.debug("proposer sent a grpc fetch request to "+strconv.Itoa(int(client_r.name)), 0)

		if err == nil && resp != nil {
			//prop.debug("proposer received a grpc fetch response from "+strconv.Itoa(int(client_r.name)), 0)
			for i := 0; i < len(resp.ClientBatches); i++ {
				foundBtch := false
				for j := 0; j < len(cltBatches); j++ {
					if resp.ClientBatches[i].Id == cltBatches[j].Id {
						foundBtch = true
						break
					}
				}
				if !foundBtch {
					cltBatches = append(cltBatches, resp.ClientBatches[i])
				}
			}

			if len(cltBatches) == numBtches {
				found = true
				//prop.debug("proposer received all the client batches for the fetch request ", 1)
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
	return rA
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
	return rA
}

// convert to proto types

func (prop *Proposer) extractDecidedSlots(indexes []int, decisions [][]string, proposers []int32) []*ProposerMessage_DecidedSlot {
	if len(indexes) != len(decisions) {
		panic("should not happen")
	}
	rA := make([]*ProposerMessage_DecidedSlot, 0)
	for i := 0; i < len(indexes); i++ {
		rA = append(rA, &ProposerMessage_DecidedSlot{
			Index:    int64(indexes[i]),
			Ids:      decisions[i],
			Proposer: int64(proposers[i]),
		})
	}

	return rA
}

// checks if element "set" of ele1 is greater than ele2

func (prop *Proposer) isGreaterThan(ele1 RecorderResponse, ele2 *RecorderResponse_Proposal, set string) bool {
	if set == "F" {
		if ele1.F.Priority > ele2.Priority {
			return true
		}
		if ele1.F.Priority == ele2.Priority && ele1.F.ProposerId > ele2.ProposerId {
			return true
		}
		if ele1.F.Priority == ele2.Priority && ele1.F.ProposerId == ele2.ProposerId && ele1.F.ThreadId > ele2.ThreadId {
			return true
		}
		return false
	} else if set == "M" {
		if ele1.M.Priority > ele2.Priority {
			return true
		}
		if ele1.M.Priority == ele2.Priority && ele1.M.ProposerId > ele2.ProposerId {
			return true
		}
		if ele1.M.Priority == ele2.Priority && ele1.M.ProposerId == ele2.ProposerId && ele1.M.ThreadId > ele2.ThreadId {
			return true
		}
		return false
	} else {
		panic("should not happen")
	}

}

// return the maximum of "set" from all replies in R

func (prop *Proposer) getMaxFromResponses(array []RecorderResponse, set string) ProposerMessage_Proposal {
	if len(array) == 0 {
		panic("should not happen")
	}

	var max *RecorderResponse_Proposal

	if set == "F" {
		max = array[0].F
	} else if set == "M" {
		max = array[0].M
	}

	for i := 1; i < len(array); i++ {
		if prop.isGreaterThan(array[i], max, set) {
			if set == "F" {
				max = array[i].F
			} else if set == "M" {
				max = array[i].M
			}
		}
	}
	return ProposerMessage_Proposal{
		Priority:   max.Priority,
		ProposerId: max.ProposerId,
		ThreadId:   max.ThreadId,
		Ids:        max.Ids,
	}
}

// compare a client batch

func (prop *Proposer) isEqual(batch1 *ProposerMessage_ClientBatch, batch2 *ProposerMessage_ClientBatch) bool {
	if batch1.Sender != batch2.Sender {
		return false
	}
	if batch1.Id != batch2.Id {
		return false
	}
	return true
}

// compare two arrays of client batches,

func (prop *Proposer) isEqualClientBatches(batch1 []*ProposerMessage_ClientBatch, batch2 []*ProposerMessage_ClientBatch) bool {
	if len(batch1) != len(batch2) {
		return false
	}
	for i := 0; i < len(batch1); i++ {
		if !prop.isEqual(batch1[i], batch2[i]) {
			return false
		}
	}
	return true
}

//  compare two proposals

func (prop *Proposer) isEqualProposal(p ProposerMessage_Proposal, m ProposerMessage_Proposal) bool {
	if p.Priority != m.Priority {
		return false
	}
	if p.ProposerId != m.ProposerId {
		return false
	}
	if p.ThreadId != m.ThreadId {
		return false
	}
	if len(p.ClientBatches) != len(m.ClientBatches) {
		return false
	}
	if !prop.isEqualClientBatches(p.ClientBatches, m.ClientBatches) {
		return false
	}

	return true
}

// main proposer logic

func (prop *Proposer) handleProposeRequest(message ProposeRequest) ProposeResponse {
	//prop.debug("proposer received propose request from the proxy ", 0)

	S := 1*4 + 0
	P := ProposerMessage_Proposal{
		Priority:      int64(prop.hi),
		ProposerId:    prop.name,
		ThreadId:      prop.threadId,
		Ids:           message.proposalStr,
		ClientBatches: prop.getProposeClientBatches(message.proposalBtch),
	}

	//prop.debug("proposer created initial proposal ", 0)

	decidedSlots := prop.extractDecidedSlots(message.lastDecidedIndexes, message.lastDecidedDecisions, message.lastDecidedProposers)

	//prop.debug("proposer proposing for instance "+strconv.Itoa(int(message.instance)), 0)

	for true {

		Pi := make([]ProposerMessage_Proposal, prop.numReplicas)

		for i := 0; i < prop.numReplicas; i++ {
			Pi[i] = ProposerMessage_Proposal{
				Priority:      P.Priority,
				ProposerId:    P.ProposerId,
				ThreadId:      P.ThreadId,
				Ids:           P.Ids,
				ClientBatches: P.ClientBatches,
			}
		}
		if S%4 == 0 && (S > 4 || !message.isLeader) {
			for i := 0; i < prop.numReplicas; i++ {
				Pi[i].Priority = int64(rand.Intn(prop.hi-10)) + 1
			}
		}

		responses := make(chan *RecorderResponse, prop.numReplicas)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(50)*time.Second)
		//prop.debug("proposer sending rpc in parallel ", 0)
		wg := sync.WaitGroup{}
		for i := 0; i < prop.numReplicas; i++ {
			wg.Add(1)
			go func(p peer, pi ProposerMessage_Proposal, s int, decidedSlots []*ProposerMessage_DecidedSlot) {
				defer wg.Done()

				if s == 4 && prop.serverMode == 1 {
					newP := &ProposerMessage_Proposal{
						Priority:      pi.Priority,
						ProposerId:    pi.ProposerId,
						ThreadId:      pi.ThreadId,
						Ids:           pi.Ids,
						ClientBatches: make([]*ProposerMessage_ClientBatch, 0),
					}
					resp, err := p.client.ESP(ctx, &ProposerMessage{
						Sender:       prop.name,
						Index:        message.instance,
						P:            newP,
						S:            int64(s),
						DecidedSlots: decidedSlots,
					})

					if err != nil {
						return
					}

					if resp.ClientBatchesNotFound {
						//prop.debug("re-proposing because the recorder did not have the batches for instance  ", 20)

						newP = &ProposerMessage_Proposal{
							Priority:      pi.Priority,
							ProposerId:    pi.ProposerId,
							ThreadId:      pi.ThreadId,
							Ids:           pi.Ids,
							ClientBatches: pi.ClientBatches,
						}
						resp, err = p.client.ESP(ctx, &ProposerMessage{
							Sender: prop.name,
							Index:  message.instance,
							P:      newP,
							S:      int64(s),
						})
					}

					if resp != nil && resp.S > 0 {
						responses <- resp
						//prop.debug("proposer received a rpc response for instance "+strconv.Itoa(int(message.instance)), 0)
					}

					return

				} else {

					resp, err := p.client.ESP(ctx, &ProposerMessage{
						Sender:       prop.name,
						Index:        message.instance,
						P:            &pi,
						S:            int64(s),
						DecidedSlots: decidedSlots,
					})

					if err != nil {
						//panic(fmt.Sprintf("%v", err))
						return
					}

					if resp != nil && resp.S > 0 {
						responses <- resp
						//prop.debug("proposer received a rpc response for instance "+strconv.Itoa(int(message.instance)), 0)
					}
					return
				}

			}(prop.peers[i], Pi[i], S, decidedSlots)
		}
		decidedSlots = make([]*ProposerMessage_DecidedSlot, 0) // only the first try have the decided slots

		go func() {
			wg.Wait()
			cancel()
			close(responses)
		}()

		responsesArray := make([]RecorderResponse, 0)
		for r := range responses {
			responsesArray = append(responsesArray, *r)
			// close the channel once a majority of the replies are collected
			if len(responsesArray) == (prop.numReplicas/2)+1 {
				if len(responsesArray) != (prop.numReplicas/2)+1 {
					panic("should not happen")
				}
				break
			}
		}

		if len(responsesArray) < (prop.numReplicas/2)+1 {
			panic("should this happen?")
		}

		// If all replies in R have S’ = S and S%4 = 0
		allRepliesHaveS := true
		for i := 0; i < len(responsesArray); i++ {
			if responsesArray[i].S != int64(S) {
				allRepliesHaveS = false
				break
			}
		}

		if allRepliesHaveS {
			//prop.debug("thread id "+strconv.Itoa(int(prop.threadId))+" proposer received recorder responses with all same S with my S "+" for index "+strconv.Itoa(int(message.instance)), 0)
		}

		if allRepliesHaveS && S%4 == 0 { //propose phase
			//prop.debug("thread id "+fmt.Sprintf(" %v ", prop.threadId)+"proposer is processing S%4==0 responses for index "+fmt.Sprintf("%v", message.instance), 2)
			allRepliesHaveFHiFit := true
			for i := 0; i < len(responsesArray); i++ {
				if responsesArray[i].F.Priority != int64(prop.hi) {
					allRepliesHaveFHiFit = false
					break
				}
			}
			if allRepliesHaveFHiFit {
				//prop.debug("thread id "+strconv.Itoa(int(prop.threadId))+" proposer succeeded propose phase fast path for index "+strconv.Itoa(int(message.instance)), 2)
				return ProposeResponse{
					index:     int(message.instance),
					decisions: responsesArray[0].F.Ids,
					proposer:  int32(responsesArray[0].F.ProposerId),
					s:         S,
				}
			}

			// P ← maximum of F’ from all replies in R
			P = prop.getMaxFromResponses(responsesArray, "F")
			//prop.debug("thread id "+strconv.Itoa(int(prop.threadId))+"proposer did not succeed in the fast path propose phase, updated P", 2)
		} else if allRepliesHaveS && S%4 == 2 {
			//prop.debug("thread id "+strconv.Itoa(int(prop.threadId))+"proposer is processing S%4==2 responses for index "+strconv.Itoa(int(message.instance)), 0)
			maxM := prop.getMaxFromResponses(responsesArray, "M")
			if prop.isEqualProposal(P, maxM) {
				//prop.debug("thread id "+strconv.Itoa(int(prop.threadId))+" proposer succeeded  in the s %4 == 2 slow path with Max m "+" for index "+strconv.Itoa(int(message.instance)), 2)
				return ProposeResponse{
					index:     int(message.instance),
					decisions: P.Ids,
					proposer:  int32(P.ProposerId),
					s:         S,
				}
			}
		} else if allRepliesHaveS && S%4 == 3 {
			//prop.debug("thread id "+strconv.Itoa(int(prop.threadId))+"proposer is processing S%4==3 responses for index "+strconv.Itoa(int(message.instance)), 0)
			P = prop.getMaxFromResponses(responsesArray, "M")
		}

		if allRepliesHaveS {
			S = S + 1
			//prop.debug("thread id "+strconv.Itoa(int(prop.threadId))+" proposer updated S to "+strconv.Itoa(S)+" for index "+strconv.Itoa(int(message.instance)), 2)
		} else {
			//  if any reply in R has S’ > S: S, P ← S’, F’ from any reply with maximum S’
			for i := 0; i < len(responsesArray); i++ {
				if responsesArray[i].S > int64(S) {
					S = int(responsesArray[i].S)
					P = ProposerMessage_Proposal{
						Priority:   responsesArray[i].F.Priority,
						ProposerId: responsesArray[i].F.ProposerId,
						ThreadId:   responsesArray[i].F.ThreadId,
						Ids:        responsesArray[i].F.Ids,
					}
					//prop.debug("thread id "+strconv.Itoa(int(prop.threadId))+" proposer received a higher S, hence updated S "+" for index "+strconv.Itoa(int(message.instance)), 2)
				}
			}
		}
	}
	panic("should not happen")
}

// infinite loop listening to the server channel

func (prop *Proposer) runProposer() {
	go func() {
		for true {

			select {

			case fetchMessage := <-prop.proxyToProposerFetchChan:
				//prop.debug("proposer received fetch request", 0)
				prop.proposerToProxyFetchChan <- prop.handleFetchRequest(fetchMessage)
				//prop.debug("proposer sent back to response to proxy for the fetch request ", 0)
				break

			case proposeMessage := <-prop.proxyToProposerChan:
				//prop.debug("proposer received propose request", -1)
				response := prop.handleProposeRequest(proposeMessage)
				if response.index != -1 {
					prop.proposerToProxyChan <- response
					//prop.debug("proposer sent back to response to proxy for the propose request", 0)
				}
				break
			case decision := <-prop.proxyToProposerDecisionChan:
				//prop.debug("proposer received decision request", 9)
				prop.handleDecisionRequest(decision)
				break
			}
		}
	}()
}

// send the decisions to every recorder

func (prop *Proposer) handleDecisionRequest(decision Decision) {
	//prop.debug("proposer starting to handle a decision request ", 11)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	decidedSlots := prop.extractDecisionSlots(decision.indexes, decision.decisions, decision.proposers)
	//prop.debug("proposer sending rpc in parallel ", -1)
	wg := sync.WaitGroup{}
	for i := 0; i < prop.numReplicas; i++ {
		wg.Add(1)
		go func(p peer, slots []*Decisions_DecidedSlot) {
			defer wg.Done()
			_, _ = p.client.InformDecision(ctx, &Decisions{
				DecidedSlots: slots,
			})
			return

		}(prop.peers[i], decidedSlots)
	}
	go func() {
		wg.Wait()
		cancel()
		//prop.debug("proposer sent decisions to everyone", 0)
	}()
}

// convert between proto types
func (prop *Proposer) extractDecisionSlots(indexes []int, decisions [][]string, proposers []int32) []*Decisions_DecidedSlot {
	if len(indexes) != len(decisions) {
		panic("should not happen")
	}

	arr := make([]*Decisions_DecidedSlot, 0)

	for i := 0; i < len(indexes); i++ {
		arr = append(arr, &Decisions_DecidedSlot{
			Index:    int64(indexes[i]),
			Ids:      decisions[i],
			Proposer: int64(proposers[i]),
		})
	}
	return arr
}
