package raxos

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"raxos/proto"
	"strconv"
)

type RPCPair struct {
	code uint8
	Obj  proto.Serializable

	/*
		Message	Codes
			ClientRequestBatch 0
			ClientResponseBatch 1
			GenericConsensus 2
			MessageBlock 3
			MessageBlockRequest 4
			ClientStatusRequest 5
			ClientStatusResponse 6
	*/
}

type OutgoingRPC struct {
	rpcPair *RPCPair
	peer    int64
}

func (in *Instance) RegisterRPC(msgObj proto.Serializable, code uint8) {
	in.rpcTable[code] = &RPCPair{code, msgObj}
}

func (in *Instance) connectToReplicas() {
	var b [1]byte
	bs := b[:1]

	//connect to replicas
	for i := in.nodeName; i < in.numReplicas; i++ {
		for true {
			conn, err := net.Dial("tcp", in.replicaAddrList[i])
			if err == nil {
				in.replicaConnections[i] = conn
				in.outgoingReplicaWriters[i] = bufio.NewWriter(in.replicaConnections[i])
				in.incomingReplicaReaders[i] = bufio.NewReader(in.replicaConnections[i])
				binary.LittleEndian.PutUint16(bs, uint16(in.nodeName))
				_, err := conn.Write(bs)
				if err != nil {
					panic(err)
				}
				break
			}
		}
	}
	in.debug("Established all outgoing connections")
}

func (in *Instance) waitForConnections() {

	// waits for connections from my self + all the replicas with lower ids + from all the clients

	var b [1]byte
	bs := b[:1]
	in.debug("Listening to messages on " + in.replicaAddrList[in.nodeName])
	in.Listener, _ = net.Listen("tcp", in.replicaAddrList[in.nodeName])

	for i := int64(0); i < in.nodeName+1+in.numClients; i++ {
		conn, err := in.Listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			panic(err)
		}
		if _, err := io.ReadFull(conn, bs); err != nil {
			fmt.Println("Connection establish error:", err)
			panic(err)
		}
		id := int32(binary.LittleEndian.Uint16(bs))
		in.debug("Received incoming tcp connection from " + strconv.Itoa(int(id)))

		if int64(id) < in.numReplicas {
			if int64(id) != in.nodeName {
				in.replicaConnections[id] = conn
				in.outgoingReplicaWriters[id] = bufio.NewWriter(in.replicaConnections[id])
				in.incomingReplicaReaders[id] = bufio.NewReader(in.replicaConnections[id])
			}
		} else if int64(id) < in.numReplicas+in.numClients {
			in.clientConnections[int64(id)-in.numReplicas] = conn
			in.outgoingClientWriters[int64(id)-in.numReplicas] = bufio.NewWriter(in.clientConnections[int64(id)-in.numReplicas])
			in.incomingClientReaders[int64(id)-in.numReplicas] = bufio.NewReader(in.clientConnections[int64(id)-in.numReplicas])
		}

	}
	in.debug("Established connections from all nodes")
}

func (in *Instance) startConnectionListners() {
	for i := int64(0); i < in.numReplicas; i++ {
		go in.connectionListener(in.incomingReplicaReaders[i])
	}
	for i := int64(0); i < in.numClients; i++ {
		go in.connectionListener(in.incomingClientReaders[i])
	}
}

func (in *Instance) connectionListener(reader *bufio.Reader) {

	var msgType uint8
	var err error = nil

	for true {
		if msgType, err = reader.ReadByte(); err != nil {
			return
		}
		if rpair, present := in.rpcTable[msgType]; present {
			obj := rpair.Obj.New()
			if err = obj.Unmarshal(reader); err != nil {
				return
			}
			in.incomingChan <- &RPCPair{
				code: msgType,
				Obj:  obj,
			}
		} else {
			in.debug("Error: received unknown message type")
		}
	}
}

func (in *Instance) run() {
	go func() {
		for true {
			in.debug("Checking channel\n")

			replicaMessage := <-in.incomingChan
			//in.lock.Lock()
			in.debug("Received replica message")
			code := replicaMessage.code
			switch code {
			case in.clientRequestBatchRpc:
				clientRequestBatch := replicaMessage.Obj.(*proto.ClientRequestBatch)
				in.debug("Client request batch" + fmt.Sprintf("%#v", clientRequestBatch))
				in.handleClientRequestBatch(clientRequestBatch)
				break

			case in.clientResponseBatchRpc:
				clientResponseBatch := replicaMessage.Obj.(*proto.ClientResponseBatch)
				in.debug("Client response batch " + fmt.Sprintf("%#v", clientResponseBatch))
				in.handleClientResponseBatch(clientResponseBatch)
				break

			case in.genericConsensusRpc:
				genericConsensus := replicaMessage.Obj.(*proto.GenericConsensus)
				in.debug("Generic Consensus  " + fmt.Sprintf("%#v", genericConsensus))
				in.handleGenericConsensus(genericConsensus)
				break

			case in.messageBlockRpc:
				messageBlock := replicaMessage.Obj.(*proto.MessageBlock)
				in.debug("Message Block  " + fmt.Sprintf("%#v", messageBlock))
				in.handleMessageBlock(messageBlock)
				break

			case in.messageBlockRequestRpc:
				messageBlockRequest := replicaMessage.Obj.(*proto.MessageBlockRequest)
				in.debug("Message Block Request " + fmt.Sprintf("%#v", messageBlockRequest))
				in.handleMessageBlockRequest(messageBlockRequest)
				break

			case in.clientStatusRequestRpc:
				clientStatusRequest := replicaMessage.Obj.(*proto.ClientStatusRequest)
				in.debug("Client Status Request " + fmt.Sprintf("%#v", clientStatusRequest))
				in.handleClientStatusRequest(clientStatusRequest)
				break

			case in.clientStatusResponseRpc:
				clientStatusResponse := replicaMessage.Obj.(*proto.ClientStatusResponse)
				in.debug("Client Status Response " + fmt.Sprintf("%#v", clientStatusResponse))
				in.handleClientStatusResponse(clientStatusResponse)
				break
			}
			//in.lock.Unlock()
		}
	}()
}

func (in *Instance) internalSendMessage(peer int64, rpcPair *RPCPair) {
	code := rpcPair.code
	oriMsg := rpcPair.Obj
	var msg proto.Serializable
	msg = in.getNewCopyOfMessage(code, oriMsg)
	var w *bufio.Writer

	if peer < in.numReplicas {
		w = in.outgoingReplicaWriters[peer]
	} else if peer < in.numReplicas+in.numClients {
		w = in.outgoingClientWriters[peer-in.numReplicas]
	}

	err := w.WriteByte(code)
	if err != nil {
		return
	}
	msg.Marshal(w)
	err = w.Flush()
	if err != nil {
		return
	}
}

func (in *Instance) startOutgoingLinks() {
	for i := 0; i < numOutgoingThreads; i++ {
		go func() {
			for true {
				outgoingMessage := <-in.outgoingMessageChan
				in.internalSendMessage(outgoingMessage.peer, outgoingMessage.rpcPair)
			}
		}()
	}
}

func (in *Instance) sendMessage(peer int64, rpcPair RPCPair) {
	in.outgoingMessageChan <- &OutgoingRPC{
		rpcPair: &rpcPair,
		peer:    peer,
	}
}
