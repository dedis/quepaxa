package raxos

import "raxos/proto"

// handler for new client batches

func (pr Proxy) ClientBatch(batch *proto.ClientBatch) {
	// put the client batch to the store
	pr.clientBatchStore.Add(batch)
	// add the batch id to the toBeProposed array
	pr.toBeProposed = append( pr.toBeProposed, batch.Id)
	
	if len(pr.toBeProposed)> pr.batchSize{
		if pr.lastProposedIndex-pr.committedIndex < pr.pipelineLength{
			// send a new proposal Request to the Proposers
			strProposals := pr.toBeProposed
			btchProposals := make([]*proto.ClientBatch, 0)
			
			for i:=0; i<len(strProposals); i++{
				btch, ok := pr.clientBatchStore.Get(strProposals[i])
				if !ok{
					panic("batch not found for the id")
				}
				btchProposals = append(btchProposals, btch)
			}
			
			proposeIndex := pr.lastProposedIndex+1
			
			newProposalRequest := ProposeRequest{
				instance:     proposeIndex,
				proposalStr:  strProposals,
				proposalBtch: btchProposals,
				msWait:       pr.getLeaderWait(int(proposeIndex)),
			}
			
			pr.proxyToProposerChan <- newProposalRequest
			
			pr.lastProposedIndex++
			pr.toBeProposed = make([]string, 0)
		}
	}
}

// handler for client status request

func (pr Proxy) handleClientStatus(status *proto.ClientStatus) {
	if status.Operation == 1 {
		if pr.serverStarted == false {
			// initiate gRPC connections
			pr.server.StartgRPC()
			pr.serverStarted = true
		}
	}
	if status.Operation == 2 {
		// print consensus log
		pr.printConsensusLog()
	}
}
