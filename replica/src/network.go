package raxos

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"raxos/proto"
	"strconv"
	"time"
)

type RPCPair struct {
	code uint8
	Obj  proto.Serializable

	/*
		Message	Code


	*/
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
	var b [1]byte
	bs := b[:1]
	in.debug("Listening to messages on " + in.replicaAddrList[in.nodeName])
	in.Listener, _ = net.Listen("tcp", in.replicaAddrList[in.nodeName])

	for i := int64(0); i < in.numReplicas+in.numClients; i++ {
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
			in.replicaConnections[id] = conn
			in.outgoingReplicaWriters[id] = bufio.NewWriter(in.replicaConnections[id])
			in.incomingReplicaReaders[id] = bufio.NewReader(in.replicaConnections[id])
		} else if int64(id) < in.numReplicas+in.numClients {
			in.clientConnections[int64(id)-in.numReplicas] = conn
			in.outgoingClientWriters[int64(id)-in.numReplicas] = bufio.NewWriter(in.clientConnections[int64(id)-in.numReplicas])
			in.incomingClientReaders[int64(id)-in.numReplicas] = bufio.NewReader(in.clientConnections[int64(id)-in.numReplicas])
		}

		// FROM HERE

		go in.replicaListener(in.incomingPeerReaders[id])

	}
	in.debug("Established connections from all nodes")
}

func (in *Instance) replicaListener(reader *bufio.Reader) {
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
			in.replicaChan <- replicaMessage{
				obj:  obj,
				code: msgType,
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

			replicaMessage := <-in.replicaChan
			in.lock.Lock()

			in.debug("Received replica message")
			code := replicaMessage.code
			switch code {
			case in.prepareRpc:
				prepare := replicaMessage.obj.(*mencius.Prepare)
				//got a Prepare message
				in.debug("Prepare" + fmt.Sprintf("%#v", prepare))
				in.handlePrepare(prepare)
				break

			case in.ackRpc:
				ack := replicaMessage.obj.(*mencius.Ack)
				//got an ack message
				in.debug("Ack " + fmt.Sprintf("%#v", ack))
				in.handleAck(ack)
				break

			case in.proposeRpc:
				propose := replicaMessage.obj.(*mencius.Propose)
				//got a Propose message
				in.debug("Propose " + fmt.Sprintf("%#v", propose))
				in.handlePropose(propose)
				break

			case in.acceptRpc:
				accept := replicaMessage.obj.(*mencius.Accept)
				//got an Accept message
				in.debug("Accept " + fmt.Sprintf("%#v", accept))
				in.handleAccept(accept)
				break

			case in.learnRpc:
				learn := replicaMessage.obj.(*mencius.Learn)
				//got a learn message
				in.debug("Learn " + fmt.Sprintf("%#v", learn))
				in.handleLearn(learn)
				break

			case in.emptyRpc:
				empty := replicaMessage.obj.(*mencius.Empty)
				//got a empty
				in.debug("Empty " + fmt.Sprintf("%#v", empty))
				in.handleEmpty(empty)
				break

			}
			in.lock.Unlock()
		}
	}()
}

func (in *Instance) unicastMsg(peer int, code uint8, msg mencius.Serializable) {
	in.outgoingMessageChannels[peer] <- RPCPair{
		code: code,
		Obj:  msg,
	}
}

func (in *Instance) broadcastMsg(code uint8, msg mencius.Serializable) {
	for peer := int64(0); peer < in.numNodes; peer++ {

		switch code {
		case in.prepareRpc:
			prepare := msg.(*mencius.Prepare)
			in.outgoingMessageChannels[peer] <- RPCPair{
				code: code,
				Obj:  prepare,
			}
			break
		case in.ackRpc:
			ack := msg.(*mencius.Ack)
			in.outgoingMessageChannels[peer] <- RPCPair{
				code: code,
				Obj:  ack,
			}
			break
		case in.proposeRpc:
			propose := msg.(*mencius.Propose)
			in.outgoingMessageChannels[peer] <- RPCPair{
				code: code,
				Obj:  propose,
			}
			break
		case in.acceptRpc:
			accept := msg.(*mencius.Accept)
			in.outgoingMessageChannels[peer] <- RPCPair{
				code: code,
				Obj:  accept,
			}
			break
		case in.learnRpc:
			learn := msg.(*mencius.Learn)
			in.outgoingMessageChannels[peer] <- RPCPair{
				code: code,
				Obj:  learn,
			}
			break
		case in.emptyRpc:
			empty := msg.(*mencius.Empty)
			in.outgoingMessageChannels[peer] <- RPCPair{
				code: code,
				Obj:  empty,
			}
			break
		}

	}
}

func (in *Instance) startHeartBeats() {
	in.debug("Sending heart beats\n")

	go func() {
		for true {
			time.Sleep(time.Duration(in.keepAliveTimeout/2) * time.Millisecond)
			empty := mencius.Empty{
				NodeId: in.nodeId,
			}

			in.debug("Heart beat from self")
			in.broadcastMsg(in.emptyRpc, &empty)
		}
	}()

}

func (in *Instance) handleEmpty(empty *mencius.Empty) {
	node := empty.NodeId
	in.lastSeenTime[node] = time.Now()
}

func (in *Instance) startFailureDetector() {
	go func() {
		in.debug("Starting failure detector")
		for true {

			time.Sleep(time.Duration(in.keepAliveTimeout) * time.Millisecond)
			if time.Now().Sub(in.lastSeenTime[in.nodeId]).Milliseconds() > in.keepAliveTimeout {
				continue
			}

			for i := int64(0); i < in.numNodes; i++ {
				if i == in.nodeId {
					continue
				}
				if time.Now().Sub(in.lastSeenTime[i]).Milliseconds() > in.keepAliveTimeout {
					// node i have failed
					in.lock.Lock()
					in.onSuspect(i)
					in.lock.Unlock()

				}
			}
		}
	}()
}

func (in *Instance) startOutgoingLinks() {

	for i := int64(0); i < in.numNodes; i++ {
		go func(peerId int64) {
			w := in.outgoingPeerWriters[peerId]

			for true {
				rpcPair := <-in.outgoingMessageChannels[peerId]
				code := rpcPair.code
				msg := rpcPair.Obj

				err := w.WriteByte(code)
				if err != nil {
					continue
				}
				msg.Marshal(w)
				err = w.Flush()
				if err != nil {
					continue
				}
			}
		}(i)
	}

}
