package raxos

import (
	"raxos/proto"
	"strconv"
)

// handler for new client batches

func (pr *Proxy) handleClientBatch(batch *proto.ClientBatch) {
	// put the client batch to the store
	pr.clientBatchStore.Add(batch)
	// add the batch id to the toBeProposed array
	pr.toBeProposed = append(pr.toBeProposed, batch.Id)

	if len(pr.toBeProposed) > pr.batchSize {
		if pr.lastProposedIndex-pr.committedIndex < pr.pipelineLength {
			// send a new proposal Request to the Proposers
			strProposals := pr.toBeProposed
			btchProposals := make([]*proto.ClientBatch, 0)

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
				uniqueID:             strconv.Itoa(uniqueId) + "." + strconv.Itoa(int(pr.name)),
				lastDecidedIndexes:   pr.lastDecidedIndexes,
				lastDecidedDecisions: pr.lastDecidedDecisions,
				lastDecidedUniqueIds: pr.lastDecidedUniqueIds,
			}

			pr.proxyToProposerChan <- newProposalRequest
			pr.replicatedLog[proposeIndex] = Slot{
				proposedBatch: strProposals,
				decidedBatch:  nil,
				uniqueId:      strconv.Itoa(uniqueId) + "." + strconv.Itoa(int(pr.name)),
			}

			pr.lastProposedIndex++
			pr.toBeProposed = make([]string, 0)
			pr.proposalId++
			pr.lastDecidedIndexes = make([]int, 0)
			pr.lastDecidedDecisions = make([][]string, 0)
			pr.lastDecidedUniqueIds = make([]string, 0)
		}
	}
}

// handler for client status request

func (pr *Proxy) handleClientStatus(status *proto.ClientStatus) {
	if status.Operation == 1 {
		if pr.serverStarted == false {
			// initiate gRPC connections
			pr.server.StartProposers()
			pr.serverStarted = true
		}
	}
	if status.Operation == 2 {
		// print consensus log
		pr.printConsensusLog()
	}
}
