package raxos

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

// start listening to the proxy tcp connection, and setup all outgoing setting (without making connections)

func (pr *Proxy) NetworkInit() {
	go pr.WaitForConnections()
	pr.StartOutgoingLinks()
}

/*
	Fill the RPC table by assigning a unique id to each message type
*/

func (pr *Proxy) RegisterRPC(msgObj proto.Serializable, code uint8) {
	pr.rpcTable[code] = &common.RPCPair{code, msgObj}
}

/*
	Listen on the server port for new connections
	Each proxy receives connection from all clients
*/

func (pr *Proxy) WaitForConnections() {

	var b [4]byte
	bs := b[:4]
	pr.debug("Listening to messages on "+pr.serverAddress, 0)
	pr.Listener, _ = net.Listen("tcp", pr.serverAddress)

	for true {
		conn, err := pr.Listener.Accept()
		if err != nil {
			fmt.Println("TCP accept error:", err)
			panic(err)
		}
		if _, err := io.ReadFull(conn, bs); err != nil {
			fmt.Println("Connection id reading error:", err)
			panic(err)
		}
		id := int32(binary.LittleEndian.Uint16(bs))
		pr.debug("Received incoming tcp connection from "+strconv.Itoa(int(id)), 0)

		pr.incomingClientReaders[int64(id)] = bufio.NewReader(conn)
		go pr.connectionListener(pr.incomingClientReaders[int64(id)], id)
		pr.debug("Started listening to "+strconv.Itoa(int(id)), 0)
		pr.connectToClient(id) // make a TCP connection with client id
	}
}

/*
	listen to a given connection. Upon receiving any message, put it into the central buffer
*/

func (pr *Proxy) connectionListener(reader *bufio.Reader, id int32) {

	var msgType uint8
	var err error = nil

	for true {
		if msgType, err = reader.ReadByte(); err != nil {
			pr.debug("Error while reading code byte: the TCP connection was broken for "+strconv.Itoa(int(id)), 0)
			return
		}
		if rpair, present := pr.rpcTable[msgType]; present {
			obj := rpair.Obj.New()
			if err = obj.Unmarshal(reader); err != nil {
				pr.debug("Error while unmarshalling", 0)
				return
			}
			pr.incomingChan <- common.RPCPair{
				Code: msgType,
				Obj:  obj,
			}
		} else {
			pr.debug("Error: received unknown message type", 0)
		}
	}
}

/*
	Make a TCP connection to the client id
*/

func (pr *Proxy) connectToClient(id int32) {
	var b [4]byte
	bs := b[:4]
	for true {
		conn, err := net.Dial("tcp", pr.clientAddrList[int64(id)])
		if err == nil {
			pr.outgoingClientWriters[int64(id)] = bufio.NewWriter(conn)
			binary.LittleEndian.PutUint16(bs, uint16(pr.name))
			_, err := conn.Write(bs)
			if err != nil {
				pr.debug("Error connecting to client "+strconv.Itoa(int(id)), 0)
				panic(err)
			}
			pr.debug("Started outgoing connection to client"+strconv.Itoa(int(id)), 0)
			break
		}
	}

}

/*
	Write a message to the wire, first the message type is written and then the actual message
*/

func (pr *Proxy) internalSendMessage(peer int64, rpcPair *common.RPCPair) {
	code := rpcPair.Code
	msg := rpcPair.Obj

	var w *bufio.Writer

	w = pr.outgoingClientWriters[peer]
	
	pr.buffioWriterMutexes[peer].Lock()

	err := w.WriteByte(code)
	if err != nil {
		pr.debug("Error while writing byte", 0)
		return
	}
	err = msg.Marshal(w)
	if err != nil {
		pr.debug("Error while marshalling", 0)
		return
	}
	err = w.Flush()
	if err != nil {
		pr.debug("Error while flushing", 0)
		return
	}
	pr.buffioWriterMutexes[peer].Unlock()
}

/*
	A set of threads that manages outgoing messages: write the message to the OS buffers
*/

func (pr *Proxy) StartOutgoingLinks() {
	for i := 0; i < 200; i++ {
		go func() {
			for true {
				outgoingMessage := <-pr.outgoingMessageChan
				pr.internalSendMessage(outgoingMessage.Peer, outgoingMessage.RpcPair)
			}
		}()
	}
}

/*
	adds a new out going message to the out going channel
    note that the message object inside the RPC pair should be unique because the protobuf objects are not thread safe
*/

func (pr *Proxy) sendMessage(peer int64, rpcPair common.RPCPair) {
	pr.outgoingMessageChan <- common.OutgoingRPC{
		RpcPair: &rpcPair,
		Peer:    peer,
	}
}