package cmd

import (
	"fmt"
	"raxos/proto"
	raxos "raxos/replica/src"
	"time"
)

func (cl *Client) handleClientStatusResponse(response *proto.ClientStatusResponse) {
	fmt.Printf("Status response from %d, with message %s\n", response.Sender, response.Message)
}

/*
	Send a status request to all the replicas
*/

func (cl *Client) SendStatus(operationType int64) {
	for i := int64(0); i < cl.numReplicas; i++ {
		statusRequest := proto.ClientStatusRequest{
			Sender:    cl.clientName,
			Receiver:  i,
			Operation: operationType,
			Message:   "",
		}

		rpcPair := raxos.RPCPair{
			Code: cl.clientStatusRequestRpc,
			Obj:  &statusRequest,
		}

		cl.sendMessage(i, rpcPair)
	}
	time.Sleep(time.Duration(cl.replicaTimeout) * time.Second)
}
