package cmd

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"raxos/common"
	"raxos/proto"
	"strconv"
)

/*
	Fill the RPC table by assigning a unique id to each message type
*/

func (cl *Client) RegisterRPC(msgObj proto.Serializable, code uint8) {
	cl.rpcTable[code] = &common.RPCPair{Code: code, Obj: msgObj}
}

/*
	Each client sends connection requests to all replica proxies
*/

func (cl *Client) ConnectToReplicas() {

	cl.debug("Connecting to "+strconv.Itoa(int(cl.numReplicas))+" replicas", 0)

	var b [4]byte
	bs := b[:4]

	for i, _ := range cl.replicaAddrList {
		for true {
			conn, err := net.Dial("tcp", cl.replicaAddrList[i])
			if err == nil {
				cl.outgoingReplicaWriters[i] = bufio.NewWriter(conn)
				binary.LittleEndian.PutUint16(bs, uint16(cl.clientName))
				_, err := conn.Write(bs)
				if err != nil {
					cl.debug("Error while writing name to replica"+strconv.Itoa(int(i)), 0)
					panic(err)
				}
				cl.debug("Established outgoing connection to "+strconv.Itoa(int(i)), 0)
				break
			}
		}
	}
	cl.debug("Established all outgoing connections to replicas", 0)
}

/*
	Listen on the port for new connections
	Whenever the proxy receives a new client connection, it dials the client
*/

func (cl *Client) WaitForConnections() {

	var b [4]byte
	bs := b[:4]
	Listener, _ := net.Listen("tcp", cl.clientListenAddress)
	cl.debug("Listening to incoming connections on "+cl.clientListenAddress, 0)

	for true {
		conn, err := Listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			panic(err)
		}
		if _, err := io.ReadFull(conn, bs); err != nil {
			fmt.Println("Connection establish error:", err)
			panic(err)
		}
		id := int32(binary.LittleEndian.Uint16(bs))
		cl.debug("Received incoming tcp connection from "+strconv.Itoa(int(id)), 0)

		cl.incomingReplicaReaders[int64(id)] = bufio.NewReader(conn)
		go cl.connectionListener(cl.incomingReplicaReaders[int64(id)], id)
		cl.debug("Started listening to "+strconv.Itoa(int(id)), 0)

	}
}

/*
	listen to a given connection. Upon receiving any message, put it into the central buffer
*/

func (cl *Client) connectionListener(reader *bufio.Reader, id int32) {

	var msgType uint8
	var err error = nil

	for true {
		if msgType, err = reader.ReadByte(); err != nil {
			cl.debug("Error while reading message code: connection broken from "+strconv.Itoa(int(id)), 1)
			return
		}
		if rpair, present := cl.rpcTable[msgType]; present {
			obj := rpair.Obj.New()
			if err = obj.Unmarshal(reader); err != nil {
				cl.debug("Error while unmarshalling", 1)
				return
			}
			cl.incomingChan <- &common.RPCPair{
				Code: msgType,
				Obj:  obj,
			}
		} else {
			cl.debug("Error: received unknown message type", 1)
		}
	}
}

/*
	This is an execution thread that listens to all the incoming messages
	It listens to incoming messages from the incomingChan, and invokes the appropriate handler depending on the message type
*/

func (cl *Client) Run() {
	go func() {
		for true {

			cl.debug("Checking channel\n", 0)
			replicaMessage := <-cl.incomingChan
			cl.debug("Received replica message", 0)
			code := replicaMessage.Code
			switch code {

			case cl.clientBatchRpc:
				clientResponseBatch := replicaMessage.Obj.(*proto.ClientBatch)
				cl.debug("Client response batch from "+fmt.Sprintf("%#v", clientResponseBatch.Sender), 0)
				cl.handleClientResponseBatch(clientResponseBatch)
				break

			case cl.clientStatusRpc:
				clientStatusResponse := replicaMessage.Obj.(*proto.ClientStatus)
				cl.debug("Client Status Response from"+fmt.Sprintf("%#v", clientStatusResponse.Sender), 0)
				cl.handleClientStatusResponse(clientStatusResponse)
				break
			}
		}
	}()
}

/*
	Write a message to the wire, first the message type is written and then the actual message
*/

func (cl *Client) internalSendMessage(peer int64, rpcPair *common.RPCPair) {
	code := rpcPair.Code
	msg := rpcPair.Obj
	var w *bufio.Writer

	w = cl.outgoingReplicaWriters[peer]
	cl.socketMutexs[peer].Lock()
	err := w.WriteByte(code)
	if err != nil {
		cl.debug("Error writing message code byte:"+err.Error(), 1)
		return
	}
	err = msg.Marshal(w)
	if err != nil {
		cl.debug("Error while marshalling:"+err.Error(), 1)
		return
	}
	err = w.Flush()
	if err != nil {
		cl.debug("Error while flushing:"+err.Error(), 1)
		return
	}
	cl.socketMutexs[peer].Unlock()
	cl.debug("Internal sent message to "+strconv.Itoa(int(peer)), 0)
}

/*
	A set of threads that manages outgoing messages: write the message to the OS buffers
*/

func (cl *Client) StartOutgoingLinks() {
	for i := 0; i < numOutgoingThreads; i++ {
		go func() {
			for true {
				outgoingMessage := <-cl.outgoingMessageChan
				cl.internalSendMessage(outgoingMessage.Peer, outgoingMessage.RpcPair)
				cl.debug("Invoked internal sent to replica "+strconv.Itoa(int(outgoingMessage.Peer)), 0)
			}
		}()
	}
}

/*
	adds a new out-going message to the out going channel
*/

func (cl *Client) sendMessage(peer int64, rpcPair common.RPCPair) {
	cl.outgoingMessageChan <- &common.OutgoingRPC{
		RpcPair: &rpcPair,
		Peer:    peer,
	}
	cl.debug("Added RPC pair to outgoing channel to peer "+strconv.Itoa(int(peer)), 0)
}
