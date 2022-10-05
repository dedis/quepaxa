package raxos

import (
	"log"
	"os"
	"raxos/proto"
	"strconv"
)

// handler for new client batches

func (pr *Proxy) handleClientBatch(batch proto.ClientBatch) {
	// put the client batch to the store
	pr.clientBatchStore.Add(batch)
	// add the batch id to the toBeProposed array
	pr.toBeProposed = append(pr.toBeProposed, batch.Id)

	if len(pr.toBeProposed) > pr.batchSize { // if we have a sufficient batch size
		if pr.lastProposedIndex-pr.committedIndex < pr.pipelineLength {
			// send a new proposal Request to the ProposersChan
			strProposals := pr.toBeProposed
			btchProposals := make([]proto.ClientBatch, 0)

			for i := 0; i < len(strProposals); i++ {
				btch, ok := pr.clientBatchStore.Get(strProposals[i])
				if !ok {
					panic("batch not found for the id")
				}
				btchProposals = append(btchProposals, btch)
			}

			proposeIndex := pr.lastProposedIndex + 1
			uniqueId := pr.proposalId

			newProposalRequest := ProposeRequest{
				instance:             proposeIndex,
				proposalStr:          strProposals,
				proposalBtch:         btchProposals,
				msWait:               pr.getLeaderWait(int(proposeIndex)),
				uniqueID:             strconv.Itoa(int(pr.name)) + "." + strconv.Itoa(uniqueId),
				lastDecidedIndexes:   pr.lastDecidedIndexes,
				lastDecidedDecisions: pr.lastDecidedDecisions,
				lastDecidedUniqueIds: pr.lastDecidedUniqueIds,
			}

			pr.proxyToProposerChan <- newProposalRequest

			// create the slot index from here

			if int64(len(pr.replicatedLog)) != proposeIndex {
				panic("propose index and length of replicated log do not match")
			}

			// create the new entry
			pr.replicatedLog = append(pr.replicatedLog, Slot{})

			pr.replicatedLog[proposeIndex] = Slot{
				proposedBatch:    strProposals,
				decidedBatch:     nil,
				proposedUniqueId: newProposalRequest.uniqueID,
				decidedUniqueId:  "",
				decided:          false,
				committed:        false,
			}

			// reset the variables
			pr.toBeProposed = make([]string, 0)
			pr.lastProposedIndex++
			pr.proposalId++
			pr.lastDecidedIndexes = make([]int, 0)
			pr.lastDecidedDecisions = make([][]string, 0)
			pr.lastDecidedUniqueIds = make([]string, 0)
		}
	}
}

// handler for client status request

func (pr *Proxy) handleClientStatus(status proto.ClientStatus) {
	if status.Operation == 1 {
		if pr.serverStarted == false {
			// initiate gRPC connections
			pr.server.StartProposers()
			pr.serverStarted = true
		}
	}
	if status.Operation == 2 {
		// print logs
		pr.printLog()
	}
}

// print the mempool and the consensus log to files

func (pr *Proxy) printLog() {
	pr.clientBatchStore.printStore(pr.logFilePath, pr.name) // print mem pool
	pr.printConsensusLog()                                  // print the replicated log
}

// print the replicated log to a file, only client batch ids are printed

func (pr *Proxy) printConsensusLog() {
	f, err := os.Create(pr.logFilePath + strconv.Itoa(int(pr.name)) + "-consensus.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for i := 0; i < len(pr.replicatedLog); i++ {
		if pr.replicatedLog[i].decided == true {
			for j := 0; j < len(pr.replicatedLog[i].decidedBatch); j++ {
				_, _ = f.WriteString(strconv.Itoa(i) + "." + strconv.Itoa(j) + ":" + pr.replicatedLog[i].decidedBatch[j] + "\n")
			}
		} else {
			break
		}
	}
}
