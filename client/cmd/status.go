package cmd

import (
	"fmt"
	"raxos/proto"
	raxos "raxos/replica/src"
	"time"
)

/*
	When a status response is received, print it to console
*/

func (cl *Client) handleClientStatusResponse(response *proto.ClientStatusResponse) {
	fmt.Printf("Status response from %d, with message %s\n", response.Sender, response.Message)
}

/*
	Send a status request to all the replicas
*/

func (cl *Client) SendStatus(operationType int64) {
	for i := int64(0); i < cl.numReplicas; i++ {

		/*
			Since the send Message doesn't expect broadcast, create a message for each replica
		*/

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
	time.Sleep(time.Duration(statusTimeout) * time.Second)
}
