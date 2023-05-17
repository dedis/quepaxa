package cmd

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"raxos/common"
	"raxos/proto"
	"raxos/proto/client"
	"strconv"
)

// this file defines the tcp methods used by client to send and receive messages from replicas

/*
	fill the RPC table by assigning a unique id to each message type
*/

func (cl *Client) RegisterRPC(msgObj proto.Serializable, code uint8) {
	cl.rpcTable[code] = &common.RPCPair{Code: code, Obj: msgObj}
}

/*
	each client sends connection requests to all replica proxies
*/

func (cl *Client) ConnectToReplicas() {

	cl.debug("Connecting to "+strconv.Itoa(int(cl.numReplicas))+" replicas", 0)

	var b [4]byte
	bs := b[:4]

	for i, _ := range cl.replicaAddrList {
		for true {
			conn, err := net.Dial("tcp", cl.replicaAddrList[i])
			if err == nil {
				// save the connection info
				cl.outgoingReplicaWriters[i] = bufio.NewWriter(conn)
				// send my name to the replica
				binary.LittleEndian.PutUint16(bs, uint16(cl.name))
				_, err := conn.Write(bs)
				if err != nil {
					cl.debug("Error while writing name to replica"+strconv.Itoa(int(i)), 0)
					panic(err)
				}
				cl.debug("established outgoing connection to "+strconv.Itoa(int(i)), -1)
				break
			} else {
				panic(err.Error())
			}
		}
	}
	cl.debug("Established all outgoing connections to replicas", 0)
}

/*
	tcp listen on the port for new connections from replicas
*/

func (cl *Client) WaitForConnections() {

	var b [4]byte
	bs := b[:4]
	Listener, err := net.Listen("tcp", cl.clientListenAddress)
	if err != nil {
		panic(err.Error())
	}
	cl.debug("Listening to incoming connections on "+cl.clientListenAddress, 0)

	for true {
		conn, err := Listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			panic(err)
		}
		// read the proxy id
		if _, err := io.ReadFull(conn, bs); err != nil {
			fmt.Println("Connection establish error:", err)
			panic(err)
		}
		id := int32(binary.LittleEndian.Uint16(bs))
		cl.debug("Received incoming tcp connection from "+strconv.Itoa(int(id)), -1)
		// save the id and the connection reader
		cl.incomingReplicaReaders[int64(id)] = bufio.NewReader(conn)
		// asynchronously listen to the messages from the new connection
		go cl.connectionListener(cl.incomingReplicaReaders[int64(id)], id)
		cl.debug("Started listening to "+strconv.Itoa(int(id)), 0)

	}
}

/*
	listen to a given tcp connection. Upon receiving any message, put it into the incoming chanel
	first the message type is received as a single byte, then the actual message is sent
*/

func (cl *Client) connectionListener(reader *bufio.Reader, id int32) {

	var msgType uint8
	var err error = nil

	for true {
		if msgType, err = reader.ReadByte(); err != nil {
			cl.debug("Error while reading message code: connection broken from "+strconv.Itoa(int(id)), 1)
			return
		}
		// proceed only if this message type is present
		if rpair, present := cl.rpcTable[msgType]; present {
			obj := rpair.Obj.New()
			if err = obj.Unmarshal(reader); err != nil {
				cl.debug("Error while unmarshalling", 1)
				return
			}
			// put a new message to the incoming Chan
			cl.incomingChan <- &common.RPCPair{
				Code: msgType,
				Obj:  obj,
			}
		} else {
			cl.debug("Error: received unknown message type", 1)
			return
		}
	}
}

/*
	execution thread that listens to all the incoming messages
	It listens to incoming messages from the incomingChan, and invokes the appropriate handler depending on the message type
*/

func (cl *Client) Run() {
	go func() {
		for true {
			cl.debug("Checking channel\n", -1)
			replicaMessage := <-cl.incomingChan
			cl.debug("Received message", 0)
			code := replicaMessage.Code
			switch code {

			case cl.clientBatchRpc:
				clientResponseBatch := replicaMessage.Obj.(*client.ClientBatch)
				cl.debug("client response batch from "+strconv.Itoa(int(clientResponseBatch.Sender)), 0)
				cl.handleClientResponseBatch(clientResponseBatch)
				break

			case cl.clientStatusRpc:
				clientStatusResponse := replicaMessage.Obj.(*client.ClientStatus)
				cl.debug("Client Status Response from "+strconv.Itoa(int(clientStatusResponse.Sender)), 0)
				cl.handleClientStatusResponse(clientStatusResponse)
				break
			}
		}
	}()
}

/*
	write a message to the wire, first the message type is written and then the actual message
*/

func (cl *Client) internalSendMessage(peer int64, rpcPair *common.RPCPair) {
	code := rpcPair.Code
	msg := rpcPair.Obj
	var w *bufio.Writer
	// mutual exclusion for writers
	w = cl.outgoingReplicaWriters[peer]
	cl.socketMutexs[peer].Lock()
	err := w.WriteByte(code)
	if err != nil {
		cl.debug("Error writing message code byte:"+err.Error(), 1)
		cl.socketMutexs[peer].Unlock()
		return
	}
	err = msg.Marshal(w)
	if err != nil {
		cl.debug("Error while marshalling:"+err.Error(), 1)
		cl.socketMutexs[peer].Unlock()
		return
	}
	err = w.Flush()
	if err != nil {
		cl.debug("Error while flushing:"+err.Error(), 1)
		cl.socketMutexs[peer].Unlock()
		return
	}
	cl.socketMutexs[peer].Unlock()
	cl.debug("internal sent message to "+strconv.Itoa(int(peer)), -1)
}

/*
	a set of threads that manages outgoing messages: write the message to the OS buffers
*/

func (cl *Client) StartOutgoingLinks() {
	for i := 0; i < numOutgoingThreads; i++ {
		go func() {
			for true {
				// get a new message, and invoke the internSend, which will write the message to buffer
				outgoingMessage := <-cl.outgoingMessageChan
				cl.internalSendMessage(outgoingMessage.Peer, outgoingMessage.RpcPair)
				cl.debug("Invoked internal sent to replica "+strconv.Itoa(int(outgoingMessage.Peer)), 0)
			}
		}()
	}
}

/*
	this is the api used by all network writes
	adds a new out-going message to the out going channel
*/

func (cl *Client) sendMessage(peer int64, rpcPair common.RPCPair) {
	cl.outgoingMessageChan <- &common.OutgoingRPC{
		RpcPair: &rpcPair,
		Peer:    peer,
	}
	cl.debug("client added RPC pair to outgoing channel to peer "+strconv.Itoa(int(peer)), -1)
}
