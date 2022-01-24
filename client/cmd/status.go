package cmd

import (
	"raxos/proto"
	raxos "raxos/replica/src"
	"strconv"
	"time"
)

/*
	When a status response is received, print it to console
*/

func (cl *Client) handleClientStatusResponse(response *proto.ClientStatusResponse) {
	cl.debug("Status response from " + strconv.Itoa(int(response.Sender)))
}

/*
	Send a status request to all the replicas
*/

func (cl *Client) SendStatus(operationType int64) {
	cl.debug("Sending status request")
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
		cl.debug("Sent status to " + strconv.Itoa(int(i)))
	}
	time.Sleep(time.Duration(statusTimeout) * time.Second)
}
