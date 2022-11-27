package cmd

import (
	"fmt"
	"raxos/common"
	"raxos/proto/client"
	"strconv"
	"time"
)

/*
	when a status response is received, print it to console
*/

func (cl *Client) handleClientStatusResponse(response *client.ClientStatus) {
	fmt.Print("Status response from " + strconv.Itoa(int(response.Sender)) + " \n")
}

/*
	Send a status request to all the replicas
*/

func (cl *Client) SendStatus(operationType int64) {

	cl.debug("Sending status request to all proxies", 0)

	for i, _ := range cl.replicaAddrList {

		statusRequest := client.ClientStatus{
			Sender:    cl.name,
			Operation: operationType,
			Message:   cl.slowdown,
		}

		rpcPair := common.RPCPair{
			Code: cl.clientStatusRpc,
			Obj:  &statusRequest,
		}

		cl.sendMessage(i, rpcPair)
		cl.debug("Sent status to "+strconv.Itoa(int(i)), 0)
	}
	time.Sleep(time.Duration(statusTimeout) * time.Second)
}
