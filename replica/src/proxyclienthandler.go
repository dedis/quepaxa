package raxos

import (
	"fmt"
	"log"
	"os"
	"raxos/common"
	"raxos/proto/client"
	"strconv"
	"time"
)

// this file defines the client handling part of Proxy QuePaxa

// handler for new client batches

func (pr *Proxy) handleClientBatch(batch client.ClientBatch) {
	// put the client batch to the store
	pr.clientBatchStore.Add(batch)
	// add the batch id to the toBeProposed array
	pr.toBeProposed = append(pr.toBeProposed, batch.Id)

	if time.Now().Sub(pr.lastTimeProposed).Microseconds() >= pr.batchTime {
		if pr.lastProposedIndex-pr.committedIndex < pr.pipelineLength {
			proposeIndex := pr.lastProposedIndex + 1
			for proposeIndex+1 <= int64(len(pr.replicatedLog)) && pr.replicatedLog[proposeIndex].decided {
				proposeIndex++ // we always propose for a new index
			}

			if pr.leaderMode == 4 && !pr.replicatedLog[proposeIndex-1].committed {
				return
			}

			msWait := int(pr.getLeaderWait(pr.getLeaderSequence(proposeIndex)))
			msWait = msWait * int(proposeIndex-pr.committedIndex) // adjust waiting for the pipelining
			if msWait < 0 {
				panic("should not happen")
			}
			if pr.instanceTimeouts[proposeIndex] != nil {
				pr.instanceTimeouts[proposeIndex].Cancel()
			}
			pr.debug("timeout for instance "+strconv.Itoa(int(proposeIndex))+"is"+strconv.Itoa(msWait), 20)
			pr.instanceTimeouts[proposeIndex] = common.NewTimerWithCancel(time.Duration(msWait) * time.Microsecond)
			pr.instanceTimeouts[proposeIndex].SetTimeoutFuntion(func() {
				pr.proposeRequestIndex <- ProposeRequestIndex{index: proposeIndex}
			})
			pr.lastProposedIndex = proposeIndex
			pr.instanceTimeouts[proposeIndex].Start()
			if msWait != 0 && (pr.leaderMode == 1 || pr.leaderMode == 2 || pr.leaderMode == 4) {
				pr.handleDecisionNotification()
			}
			pr.lastTimeProposed = time.Now()
		}
	}
}

// handler for client status request

func (pr *Proxy) handleClientStatus(status client.ClientStatus) {
	if status.Operation == 1 {
		if pr.serverStarted == false {
			// initiate gRPC connections
			pr.debug("proxy starting proposers  ", -1)
			pr.server.StartProposers()
			pr.serverStarted = true
			pr.startTime = time.Now()
		}
	}
	if status.Operation == 2 {
		pr.debug("proxy printing logs", 0)
		// print logs
		pr.printLog()
	}

	if status.Operation == 4 {
		pr.debug("printing the steps per slot", 0)
		avg, totalSlots, stepsAccum := pr.calculateStepsPerSlot()
		fmt.Printf("\nAverage number of steps per slot: %f, total slots %v, steps accumilated %v ", avg, totalSlots, stepsAccum)
	}
}

/*
	calculate the average number of steps per slot
*/

func (pr *Proxy) calculateStepsPerSlot() (float64, int, int) {
	proposedSlots := 0
	stepsAccum := 0
	for i := 0; i < len(pr.replicatedLog); i++ {
		if pr.replicatedLog[i].committed {
			if pr.replicatedLog[i].s == 0 {
				continue
			}
			proposedSlots++
			stepsAccum += pr.replicatedLog[i].s - 3
		}
	}
	if proposedSlots == 0 {
		return 0, 0, 0
	} else {
		return float64(stepsAccum) / float64(proposedSlots), proposedSlots, stepsAccum
	}

}

// print the mempool and the consensus log to files

func (pr *Proxy) printLog() {
	pr.printConsensusLog() // print the replicated log
}

// print the replicated log to a file

func (pr *Proxy) printConsensusLog() {
	f, err := os.Create(pr.logFilePath + strconv.Itoa(int(pr.name)) + "-consensus.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for i := 0; i < len(pr.replicatedLog); i++ {
		if pr.replicatedLog[i].committed == true {
			for j := 0; j < len(pr.replicatedLog[i].decidedBatch); j++ {
				batch, ok := pr.clientBatchStore.Get(pr.replicatedLog[i].decidedBatch[j])
				if !ok {
					panic("committed batch not in the store")
				} else {
					for k := 0; k < len(batch.Messages); k++ {
						_, _ = f.WriteString(strconv.Itoa(i) + "." + strconv.Itoa(j) + "." + strconv.Itoa(k) + ":" + batch.Messages[k].Message + "\n")
					}
				}
			}
		} else {
			break
		}
	}
}

// propose to index

func (pr *Proxy) proposeToIndex(proposeIndex int64) {

	if int64(len(pr.replicatedLog)) > proposeIndex && pr.replicatedLog[proposeIndex].decided == true {
		pr.debug("did not propose for index "+strconv.Itoa(int(proposeIndex))+" because it was decided", 9)
		return
	}
	pr.instanceTimeouts[proposeIndex] = nil

	pr.debug("proposing for index "+strconv.Itoa(int(proposeIndex)), 20)

	if pr.leaderMode == 2 {
		if pr.isBeginningOfEpoch(proposeIndex) {
			pr.debug("proposing the last epoch summary for index "+strconv.Itoa(int(proposeIndex))+"", 13)
			pr.proposePreviousEpochSummary(proposeIndex)
			return
		}
	}

	batchSize := pr.batchSize
	if len(pr.toBeProposed) < batchSize {
		batchSize = len(pr.toBeProposed)
	}

	strProposals := make([]string, 0)
	btchProposals := make([]client.ClientBatch, 0)

	if batchSize == 0 {
		strProposals = []string{"nil"}
		btchProposals = append(btchProposals, client.ClientBatch{
			Sender:   -1,
			Messages: nil,
			Id:       "nil",
		})
		pr.debug("proposing empty values for index "+strconv.Itoa(int(proposeIndex)), 9)
	} else {
		// send a new proposal Request to the ProposersChan
		strProposals = pr.toBeProposed[0:batchSize]
		pr.toBeProposed = pr.toBeProposed[batchSize:]
		btchProposals = make([]client.ClientBatch, 0)

		for i := 0; i < len(strProposals); i++ {
			btch, ok := pr.clientBatchStore.Get(strProposals[i])
			if !ok {
				panic("batch not found for the id")
			}
			btchProposals = append(btchProposals, btch)
		}
	}

	waitTime := int(pr.getLeaderWait(pr.getLeaderSequence(proposeIndex)))
	isLeader := true

	if pr.leaderMode == 3 {
		isLeader = false
	} else if pr.leaderMode != 3 && waitTime != 0 {
		isLeader = false
	}

	newProposalRequest := ProposeRequest{
		instance:             proposeIndex,
		proposalStr:          strProposals,
		proposalBtch:         btchProposals,
		isLeader:             isLeader,
		lastDecidedIndexes:   pr.lastDecidedIndexes,
		lastDecidedProposers: pr.lastDecidedProposers,
		lastDecidedDecisions: pr.lastDecidedDecisions,
	}

	pr.proxyToProposerChan <- newProposalRequest
	pr.debug("proxy sent a proposal request to proposer  ", -1)
	// create the slot index
	for len(pr.replicatedLog) < int(proposeIndex+1) {
		// create the new entry
		pr.replicatedLog = append(pr.replicatedLog, Slot{
			proposedBatch: nil,
			decidedBatch:  nil,
			decided:       false,
			committed:     false,
			s:             0,
		})
	}

	pr.replicatedLog[proposeIndex] = Slot{
		proposedBatch: strProposals,
		decidedBatch:  pr.replicatedLog[proposeIndex].decidedBatch,
		decided:       pr.replicatedLog[proposeIndex].decided,
		committed:     pr.replicatedLog[proposeIndex].committed,
	}

	// reset the variables
	pr.lastDecidedIndexes = make([]int, 0)
	pr.lastDecidedProposers = make([]int32, 0)
	pr.lastDecidedDecisions = make([][]string, 0)
}
