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

/*
	RPC pair assigns a unique id to each type of message defined in the proto files
*/

type RPCPair struct {
	Code uint8
	Obj  proto.Serializable
}

/*
	Outgoing RPC assigns a rpc to its intended destination peer, the peer can be a replica or a client
*/

type OutgoingRPC struct {
	RpcPair *RPCPair
	Peer    int64
}

/*
	Fill the RPC table by assigning a unique id to each message type
*/

func (in *Instance) RegisterRPC(msgObj proto.Serializable, code uint8) {
	in.rpcTable[code] = &RPCPair{code, msgObj}
}

/*
	Each replica sends connection requests to all replicas
*/

func (in *Instance) connectToReplicas() {
	var b [4]byte
	bs := b[:4]

	//connect to replicas
	for i := int64(0); i < in.numReplicas; i++ {
		for true {
			conn, err := net.Dial("tcp", in.replicaAddrList[i])
			if err == nil {
				in.outgoingReplicaWriters[i] = bufio.NewWriter(conn)
				binary.LittleEndian.PutUint16(bs, uint16(in.nodeName))
				_, err := conn.Write(bs)
				if err != nil {
					in.debug("Error making a connection to " + strconv.Itoa(int(i)))
					panic(err)
				}
				in.debug("Made outgoing connection to replica " + strconv.Itoa(int(i)))
				break
			}
		}
	}
	in.debug("Established all outgoing connections")
}

/*
	Listen on the server port for new connections
	Each replica receives connection from all replicas and from all clients
*/

func (in *Instance) WaitForConnections() {

	var b [4]byte
	bs := b[:4]
	in.debug("Listening to messages on " + in.replicaAddrList[in.nodeName])
	in.Listener, _ = net.Listen("tcp", in.replicaAddrList[in.nodeName])

	for true {
		conn, err := in.Listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			panic(err)
		}
		if _, err := io.ReadFull(conn, bs); err != nil {
			fmt.Println("Connection id reading error:", err)
			panic(err)
		}
		id := int32(binary.LittleEndian.Uint16(bs))
		in.debug("Received incoming tcp connection from " + strconv.Itoa(int(id)))

		if int64(id) < in.numReplicas {
			// the connection is from a replica
			in.incomingReplicaReaders[id] = bufio.NewReader(conn)
			go in.connectionListener(in.incomingReplicaReaders[id], id)
			in.debug("Started listening to " + strconv.Itoa(int(id)))

		} else if int64(id) < in.numReplicas+in.numClients {
			// the connection is from a client
			in.incomingClientReaders[int64(id)-in.numReplicas] = bufio.NewReader(conn)
			go in.connectionListener(in.incomingClientReaders[int64(id)-in.numReplicas], id)
			in.debug("Started listening to " + strconv.Itoa(int(id)))
			in.connectToClient(id) // make a TCP connection with client id
		}
	}
}

/*
	listen to a given connection. Upon receiving any message, put it into the central buffer
*/

func (in *Instance) connectionListener(reader *bufio.Reader, id int32) {

	var msgType uint8
	var err error = nil

	for true {
		if msgType, err = reader.ReadByte(); err != nil {
			in.debug("Error while reading code byte: the TCP connection was broken for " + strconv.Itoa(int(id)))
			return
		}
		if rpair, present := in.rpcTable[msgType]; present {
			obj := rpair.Obj.New()
			if err = obj.Unmarshal(reader); err != nil {
				in.debug("Error while unmarshalling")
				return
			}
			in.incomingChan <- &RPCPair{
				Code: msgType,
				Obj:  obj,
			}
		} else {
			in.debug("Error: received unknown message type")
		}
	}
}

/*
	This is the main execution thread
	It listens to incoming messages from the incomingChan, and invoke the appropriate handler depending on the message type
*/

func (in *Instance) Run() {
	go func() {
		for true {
			replicaMessage := <-in.incomingChan
			//in.lock.Lock()
			in.debug("Received  message")
			code := replicaMessage.Code
			switch code {

			case in.clientRequestBatchRpc:
				clientRequestBatch := replicaMessage.Obj.(*proto.ClientRequestBatch)
				in.debug("Client request batch with id " + fmt.Sprintf("%#v", clientRequestBatch.Id))
				in.handleClientRequestBatch(clientRequestBatch)
				break

			case in.genericConsensusRpc:
				genericConsensus := replicaMessage.Obj.(*proto.GenericConsensus)
				in.debug("Generic Consensus  " + fmt.Sprintf("%#v", genericConsensus))
				in.handleGenericConsensus(genericConsensus)
				break

			case in.messageBlockRpc:
				messageBlock := replicaMessage.Obj.(*proto.MessageBlock)
				in.debug("Message Block hash  " + fmt.Sprintf("%#v", messageBlock.Hash))
				in.handleMessageBlock(messageBlock)
				break

			case in.messageBlockRequestRpc:
				messageBlockRequest := replicaMessage.Obj.(*proto.MessageBlockRequest)
				in.debug("Message Block Request from " + fmt.Sprintf("%#v", messageBlockRequest.Sender))
				in.handleMessageBlockRequest(messageBlockRequest)
				break

			case in.clientStatusRequestRpc:
				clientStatusRequest := replicaMessage.Obj.(*proto.ClientStatusRequest)
				in.debug("Client Status Request from" + fmt.Sprintf("%#v", clientStatusRequest.Sender))
				in.handleClientStatusRequest(clientStatusRequest)
				break

			case in.messageBlockAckRpc:
				messageBlockAck := replicaMessage.Obj.(*proto.MessageBlockAck)
				in.debug("Message Block Ack " + fmt.Sprintf("%#v", messageBlockAck.Hash))
				in.handleMessageBlockAck(messageBlockAck)
				break

			}
			//in.lock.Unlock()
		}
	}()
}

/*
	Write a message to the wire, first the message type is written and then the actual message
*/

func (in *Instance) internalSendMessage(peer int64, rpcPair *RPCPair) {
	code := rpcPair.Code
	oriMsg := rpcPair.Obj
	var msg proto.Serializable
	msg = in.getNewCopyOfMessage(code, oriMsg)
	var w *bufio.Writer

	if peer < in.numReplicas {
		w = in.outgoingReplicaWriters[peer]
	} else if peer < in.numReplicas+in.numClients {
		w = in.outgoingClientWriters[peer-in.numReplicas]
	}

	in.buffioWriterMutexes[peer].Lock()

	err := w.WriteByte(code)
	if err != nil {
		in.debug("Error while writing byte")
		return
	}
	err = msg.Marshal(w)
	if err != nil {
		in.debug("Error while marshalling")
		return
	}
	err = w.Flush()
	if err != nil {
		in.debug("Error while flushing")
		return
	}
	in.buffioWriterMutexes[peer].Unlock()
}

/*
	A set of threads that manages outgoing messages: write the message to the OS buffers
*/

func (in *Instance) StartOutgoingLinks() {
	for i := 0; i < numOutgoingThreads; i++ {
		go func() {
			for true {
				outgoingMessage := <-in.outgoingMessageChan
				in.internalSendMessage(outgoingMessage.Peer, outgoingMessage.RpcPair)
			}
		}()
	}
}

/*
	adds a new out going message to the out going channel
*/

func (in *Instance) sendMessage(peer int64, rpcPair RPCPair) {
	in.outgoingMessageChan <- &OutgoingRPC{
		RpcPair: &rpcPair,
		Peer:    peer,
	}
}
