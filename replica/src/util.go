package raxos

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"raxos/proto"
)

func (in *Instance) debug(message string) {
	fmt.Printf("%s\n", message)

}

func getRealSizeOf(v interface{}) (int, error) {
	b := new(bytes.Buffer)
	if err := gob.NewEncoder(b).Encode(v); err != nil {
		return 0, err
	}
	return b.Len(), nil
}

func getStringOfSizeN(length int) string {
	str := "a"
	size, _ := getRealSizeOf(str)
	for size < length {
		str = str + "a"
		size, _ = getRealSizeOf(str)
	}
	return str
}

func (in *Instance) getNewCopyofMessage(code uint8, msg proto.Serializable) proto.Serializable {

	if code == in.clientRequestRpc {

		clientRequest := msg.(*proto.ClientRequest)
		return &proto.ClientRequest{
			Id:      clientRequest.Id,
			Message: clientRequest.Message,
		}

	}

	if code == in.clientResponseRpc {
		clientResponse := msg.(*proto.ClientResponse)
		return &proto.ClientResponse{
			Id:      clientResponse.Id,
			Message: clientResponse.Message,
		}

	}

	if code == in.genericConsensusRpc {

		genericConsensus := msg.(*proto.GenericConsensus)
		return &proto.GenericConsensus{
			Index:       genericConsensus.Index,
			M:           genericConsensus.M,
			S:           genericConsensus.S,
			P:           genericConsensus.P,
			E:           genericConsensus.E,
			C:           genericConsensus.C,
			D:           genericConsensus.D,
			DS:          genericConsensus.DS,
			PR:          genericConsensus.PR,
			Destination: genericConsensus.Destination,
		}

	}

	if code == in.messageBlockReplyRpc {

		messageBlockReply := msg.(*proto.MessageBlockReply)
		return &proto.MessageBlockReply{
			Hash:     messageBlockReply.Hash,
			Requests: messageBlockReply.Requests,
		}

	}

	if code == in.messageBlockRequestRpc {
		messageBlockRequest := msg.(*proto.MessageBlockRequest)
		return &proto.MessageBlockRequest{Hash: messageBlockRequest.Hash}

	}

	return nil

}
