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
	peers                    []peer // gRPC connection list (not shared)
	proxyToProposerChan      chan ProposeRequest
	proposerToProxyChan      chan ProposeResponse
	proxyToProposerFetchChan chan FetchRequest
	proposerToProxyFetchChan chan FetchResposne
	lastSeenTimes            []*time.Time // use proposer name - 1 as index
	leaderTimeout            int64
	debugOn                  bool // if turned on, the debug messages will be print on the console
	debugLevel               int  // debug level
	hi                       int  // hi priority
	serverMode               int  // if 1, use the fast path LAN optimizations
}

// instantiate a new Proposer

func NewProposer(name int64, threadId int64, peers []peer, proxyToProposerChan chan ProposeRequest, proposerToProxyChan chan ProposeResponse, proxyToProposerFetchChan chan FetchRequest, proposerToProxyFetchChan chan FetchResposne, lastSeenTimes []*time.Time, debugOn bool, debugLevel int, hi int, serverMode int, leaderTimeout int64) *Proposer {

	pr := Proposer{
		numReplicas:              len(peers),
		name:                     name,
		threadId:                 threadId + 1,
		peers:                    peers,
		proxyToProposerChan:      proxyToProposerChan,
		proposerToProxyChan:      proposerToProxyChan,
		proxyToProposerFetchChan: proxyToProposerFetchChan,
		proposerToProxyFetchChan: proposerToProxyFetchChan,
		lastSeenTimes:            lastSeenTimes,
		leaderTimeout:            leaderTimeout,
		debugOn:                  debugOn,
		debugLevel:               debugLevel,
		hi:                       hi,
		serverMode:               serverMode,
	}

	pr.debug("created a new proposer "+fmt.Sprintf("%v", pr), -1)

	return &pr
}

/*
	if turned on, print the message to console,
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
	prop.debug("proposer starting to handle a fetch request "+fmt.Sprintf("%v", message), 1)
	found := false
	cltBatches := make([]*DecideResponse_ClientBatch, 0)
	numBtches := len(message.ids)

	for !found {

		client := prop.getRandomClient()
		clientCon := client.client
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
		resp, err := clientCon.FetchBatches(ctx, &DecideRequest{
			Ids: message.ids,
		})

		prop.debug("proposer sent a grpc fetch request to "+fmt.Sprintf("%v", client), 0)

		if err == nil && resp != nil {
			prop.debug("proposer received a grpc fetch response "+fmt.Sprintf("%v", resp), 0)
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
				prop.debug("proposer received all the client batches for the fetch request "+fmt.Sprintf("%v", cltBatches), 1)
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

// return the maximum of F’ from all replies in R

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

// run the proposer logic

func (prop *Proposer) handleProposeRequest(message ProposeRequest) ProposeResponse {
	prop.debug("proposer received propose request from the proxy "+fmt.Sprintf("%v", message), -1)

	S := 1*4 + 0
	P := ProposerMessage_Proposal{
		Priority:      int64(prop.hi),
		ProposerId:    prop.name,
		ThreadId:      prop.threadId,
		Ids:           message.proposalStr,
		ClientBatches: prop.getProposeClientBatches(message.proposalBtch),
	}

	prop.debug("proposer created initial proposal "+fmt.Sprintf("%v", P)+" for index "+fmt.Sprintf("%v", message.instance), 0)

	decidedSlots := prop.extractDecidedSlots(message.lastDecidedIndexes, message.lastDecidedDecisions)

	// sleep for msWait
	time.Sleep(time.Duration(message.msWait) * time.Millisecond)

	timeStr := ""

	if prop.debugOn {

		for y := 0; y < len(prop.lastSeenTimes); y++ {
			timeStr = timeStr + fmt.Sprintf(": %v", time.Now().Sub(*prop.lastSeenTimes[y]).Milliseconds())
		}
	}

	// if there is no proposal from anyone, propose, else return

	if !prop.noProposalUntilNow() {
		prop.debug("proposer did not propose because someone else has proposed "+timeStr, 9)
		return ProposeResponse{
			index:     -1,
			decisions: nil,
		}
	}

	prop.debug("proposer proposes "+timeStr, 9)

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
				Pi[i].Priority = int64(rand.Intn(prop.hi-2)) + 1
			}

			for i := 0; i < prop.numReplicas; i++ {
				prop.debug("thread id "+fmt.Sprintf(" %v ", prop.threadId)+"+ proposal  priority for  replica "+fmt.Sprintf("%v is %v ", i, Pi[i].Priority)+" for index "+fmt.Sprintf("%v", message.instance), 2)
			}
		}

		responses := make(chan *RecorderResponse, prop.numReplicas)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100)*time.Second)
		prop.debug("proposer sending rpc in parallel ", -1)
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
					//panic(fmt.Sprintf("%v", err))
					return
				}

				if resp != nil && resp.S > 0 {
					responses <- resp
				} else {
					//panic("response :"+fmt.Sprintf("%v", resp))
				}

				prop.debug("proposer received a rpc response "+fmt.Sprintf("S: %v, F:%v, and M:%v", resp.S, resp.F, resp.M)+" for index "+fmt.Sprintf("%v", message.instance), -1)
				return

			}(prop.peers[i], Pi[i], S, decidedSlots)
		}
		decidedSlots = make([]*ProposerMessage_DecidedSlot, 0) // only the first try have the decided slots

		go func() {
			wg.Wait()
			cancel()
			close(responses)
		}()

		// todo add fast path checks where HasClientBacthes is false when s = 0, in which case the proposer sends the s=0 again with actual batches

		responsesArray := make([]RecorderResponse, 0)
		for r := range responses {
			responsesArray = append(responsesArray, *r)
			// close the channel once a majority of the replies are collected
			if len(responsesArray) == (prop.numReplicas/2)+1 {
				for i := 0; i < (prop.numReplicas/2)+1; i++ {
					prop.debug("thread id "+fmt.Sprintf(" %v ", prop.threadId)+"proposer received a rpc response "+fmt.Sprintf("S: %v, F:%v,%v,%v,%v, and M:%v,%v,%v", responsesArray[i].S, responsesArray[i].F.Priority, responsesArray[i].F.ProposerId, responsesArray[i].F.ThreadId, responsesArray[i].F.Ids[0], responsesArray[i].M.Priority, responsesArray[i].M.ProposerId, responsesArray[i].M.ThreadId)+" for index "+fmt.Sprintf("%v", message.instance), 2)
				}

				if len(responsesArray) != (prop.numReplicas/2)+1 {
					panic("should not happen")
				}
				break
			}
		}

		if len(responsesArray) == 0 {
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

		prop.debug("thread id "+fmt.Sprintf(" %v ", prop.threadId)+"proposer received recorder responses with all same S with my S "+" for index "+fmt.Sprintf("%v", message.instance), 2)

		if allRepliesHaveS && S%4 == 0 { //propose phase
			prop.debug("thread id "+fmt.Sprintf(" %v ", prop.threadId)+"proposer is processing S%4==0 responses for index "+fmt.Sprintf("%v", message.instance), 2)
			allRepliesHaveFHiFit := true
			for i := 0; i < len(responsesArray); i++ {
				if responsesArray[i].F.Priority != int64(prop.hi) {
					allRepliesHaveFHiFit = false
					break
				}
			}
			if allRepliesHaveFHiFit {
				prop.debug("thread id "+fmt.Sprintf(" %v ", prop.threadId)+"proposer succeeded propose phase fast path for index "+fmt.Sprintf("%v", message.instance), 2)
				return ProposeResponse{
					index:     int(message.instance),
					decisions: responsesArray[0].F.Ids,
				}
			}

			// P ← maximum of F’ from all replies in R
			P = prop.getMaxFromResponses(responsesArray, "F")
			prop.debug("thread id "+fmt.Sprintf(" %v ", prop.threadId)+"proposer did not succeed in the fast path propose phase, updated P to "+fmt.Sprintf("priority: %v, replica:%v, thread:%v, and ids:%v", P.Priority, P.ProposerId, P.ThreadId, P.Ids[0])+" for index "+fmt.Sprintf("%v", message.instance), 2)
		} else if allRepliesHaveS && S%4 == 2 {
			prop.debug("thread id "+fmt.Sprintf(" %v ", prop.threadId)+"proposer is processing S%4==2 responses for index "+fmt.Sprintf("%v", message.instance), 0)
			maxM := prop.getMaxFromResponses(responsesArray, "M")
			if prop.isEqualProposal(P, maxM) {
				prop.debug("thread id "+fmt.Sprintf(" %v ", prop.threadId)+"proposer succeeded  in the s %4 == 2 slow path with Max m "+fmt.Sprintf("priority: %v, proposer:%v, thread:%v, and ids:%v ,", maxM.Priority, maxM.ProposerId, maxM.ThreadId, maxM.Ids[0])+" for index "+fmt.Sprintf("%v", message.instance), 2)
				return ProposeResponse{
					index:     int(message.instance),
					decisions: P.Ids,
				}
			}
		} else if allRepliesHaveS && S%4 == 3 {
			prop.debug("thread id "+fmt.Sprintf(" %v ", prop.threadId)+"proposer is processing S%4==3 responses for index "+fmt.Sprintf("%v", message.instance), 0)
			P = prop.getMaxFromResponses(responsesArray, "M")
			prop.debug("thread id "+fmt.Sprintf(" %v ", prop.threadId)+"proposer is in S%4 ==3 gather phase and updated P to "+fmt.Sprintf("priority: %v, replica:%v, thread:%v, and ids:%v", P.Priority, P.ProposerId, P.ThreadId, P.Ids[0])+" for index "+fmt.Sprintf("%v", message.instance), 2)
		}

		if allRepliesHaveS {
			S = S + 1
			prop.debug("thread id "+fmt.Sprintf(" %v ", prop.threadId)+"proposer updated S to "+fmt.Sprintf("%v", S)+" for index "+fmt.Sprintf("%v", message.instance)+" and P to "+fmt.Sprintf("priority: %v, replica:%v, thread:%v, and ids:%v", P.Priority, P.ProposerId, P.ThreadId, P.Ids[0]), 2)
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
					prop.debug("thread id "+fmt.Sprintf(" %v ", prop.threadId)+"proposer received a higher S, hence updated S to "+fmt.Sprintf("%v", S)+" and P to "+fmt.Sprintf("priority: %v, replica:%v, thread:%v, and ids:%v", P.Priority, P.ProposerId, P.ThreadId, P.Ids[0])+" for index "+fmt.Sprintf("%v", message.instance), 2)
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
				prop.debug("proposer received fetch request", 0)
				prop.proposerToProxyFetchChan <- prop.handleFetchRequest(fetchMessage)
				prop.debug("proposer sent back to response to proxy for the fetch request ", 0)
				break

			case proposeMessage := <-prop.proxyToProposerChan:
				prop.debug("proposer received propose request", -1)
				response := prop.handleProposeRequest(proposeMessage)
				if response.index != -1 {
					prop.proposerToProxyChan <- response
					prop.debug("proposer sent back to response to proxy for the propose request to "+fmt.Sprintf("%v", response), 0)
				}
				break
			}
		}
	}()
}

// have I seen any proposal from anyone else, within the duration time.Now - leader timeout : time.Now

func (prop *Proposer) noProposalUntilNow() bool {
	for i := 0; i < len(prop.lastSeenTimes); i++ {
		if int64(i+1) == prop.name { // this hardcodes the fact that node ids start with 1
			continue
		}
		if time.Now().Sub(*prop.lastSeenTimes[i]).Milliseconds() < prop.leaderTimeout {
			return false
		}
	}

	return true
}
