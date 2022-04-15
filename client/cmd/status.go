package cmd

import (
	"fmt"
	"raxos/proto"
	raxos "raxos/replica/src"
	"strconv"
	"time"
)

/*
	When a status response is received, print it to console
*/

func (cl *Client) handleClientStatusResponse(response *proto.ClientStatusResponse) {
	fmt.Print("Status response from " + strconv.Itoa(int(response.Sender)) + " \n")
}

/*
	Send a status request to all the replicas
*/

func (cl *Client) SendStatus(operationType int64) {
	cl.debug("Sending status request to all replicas", 0)
	for i := int64(0); i < cl.numReplicas; i++ {

		/*
			Since send Message doesn't expect broadcast, create a message for each replica
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
		cl.debug("Sent status to "+strconv.Itoa(int(i)), 0)
	}
	time.Sleep(time.Duration(statusTimeout) * time.Second)
}
