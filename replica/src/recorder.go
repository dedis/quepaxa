package raxos

import (
	"math/rand"
	"raxos/proto"
	"strconv"
	"strings"
)

/*
	mark the slot as decided, and record the decided value
*/

func (in *Instance) recordRecorderDecide(consensusMessage *proto.GenericConsensus) {
	in.recorderReplicatedLog[consensusMessage.Index].decided = true
	in.recorderReplicatedLog[consensusMessage.Index].S = consensusMessage.S
	in.recorderReplicatedLog[consensusMessage.Index].decision = consensusMessage.DS
	in.recorderReplicatedLog[consensusMessage.Index].proposer = consensusMessage.PR

	in.recorderReplicatedLog[consensusMessage.Index].P = []*proto.GenericConsensusValue{}
	in.recorderReplicatedLog[consensusMessage.Index].E = []*proto.GenericConsensusValue{}
	in.recorderReplicatedLog[consensusMessage.Index].C = []*proto.GenericConsensusValue{}

	in.recorderReplicatedLog[consensusMessage.Index].proposeResponses = []*proto.GenericConsensus{}
	in.recorderReplicatedLog[consensusMessage.Index].spreadEResponses = []*proto.GenericConsensus{}
	in.recorderReplicatedLog[consensusMessage.Index].spreadCGatherEResponses = []*proto.GenericConsensus{}
	in.recorderReplicatedLog[consensusMessage.Index].gatherCResponses = []*proto.GenericConsensus{}

}

/*
	handler for the recorder. all the recorder logic goes here
*/

func (in *Instance) handleRecorderConsensusMessage(consensusMessage *proto.GenericConsensus) {

	in.recorderReplicatedLog = in.initializeSlot(in.recorderReplicatedLog, consensusMessage.Index) // create the slot if not already created

	// case 1: if this entry has been previously decided, send the response back
	if in.recorderReplicatedLog[consensusMessage.Index].decided == true && consensusMessage.D != true {
		in.debug("Received a message from "+strconv.Itoa(int(consensusMessage.Sender))+" that is not a decide message, replying with a decision message", 1)
		in.sendRecorderDecided(consensusMessage)
		return
	}

	// case 2: if this is decide message, and if I have not previously decided on this slot, decide it. This case is already handled in the common.go

	// case 3: consensus messages from a higher step
	if in.recorderReplicatedLog[consensusMessage.Index].S < consensusMessage.S {
		in.recorderReplicatedLog[consensusMessage.Index].S = consensusMessage.S

		// case 3.1 a propose message
		if consensusMessage.S%4 == 0 {
			in.handleRecorderProposeMessage(consensusMessage)
		}
	}

	// case 4 a spreadE message

	if consensusMessage.S%4 == 1 && in.recorderReplicatedLog[consensusMessage.Index].S == consensusMessage.S {
		in.handleRecorderSpreadEMessage(consensusMessage)
	}

	// case 5 a spreadCgatherE message

	if consensusMessage.S%4 == 2 && in.recorderReplicatedLog[consensusMessage.Index].S == consensusMessage.S {
		in.handleRecorderSpreadCGatherEMessage(consensusMessage)
	}

	// send the response back to proposer

	consensusReply := proto.GenericConsensus{
		Sender:      in.nodeName,
		Receiver:    consensusMessage.Sender,
		Index:       consensusMessage.Index,
		S:           in.recorderReplicatedLog[consensusMessage.Index].S,
		P:           in.removeEmptyValues(in.recorderReplicatedLog[consensusMessage.Index].P),
		E:           in.removeEmptyValues(in.recorderReplicatedLog[consensusMessage.Index].E),
		C:           in.removeEmptyValues(in.recorderReplicatedLog[consensusMessage.Index].C),
		D:           in.recorderReplicatedLog[consensusMessage.Index].decided,
		DS:          in.recorderReplicatedLog[consensusMessage.Index].decision,
		PR:          in.recorderReplicatedLog[consensusMessage.Index].proposer,
		Destination: in.consensusMessageProposerDestination,
	}

	rpcPair := RPCPair{
		Code: in.genericConsensusRpc,
		Obj:  &consensusReply,
	}
	in.debug("Recorder sending a generic consensus reply message to "+strconv.Itoa(int(consensusMessage.Sender)), 1)

	in.sendMessage(consensusMessage.Sender, rpcPair)
}

/*
	send a decide message to the proposer
*/

func (in *Instance) sendRecorderDecided(consensusMessage *proto.GenericConsensus) {
	consensusReply := proto.GenericConsensus{
		Sender:      in.nodeName,
		Receiver:    consensusMessage.Sender,
		Index:       consensusMessage.Index,
		S:           in.recorderReplicatedLog[consensusMessage.Index].S,
		P:           nil,
		E:           nil,
		C:           nil,
		D:           in.recorderReplicatedLog[consensusMessage.Index].decided,
		DS:          in.recorderReplicatedLog[consensusMessage.Index].decision,
		PR:          in.recorderReplicatedLog[consensusMessage.Index].proposer,
		Destination: in.consensusMessageCommonDestination,
	}

	rpcPair := RPCPair{
		Code: in.genericConsensusRpc,
		Obj:  &consensusReply,
	}
	in.debug("sending a decide consensus message to "+strconv.Itoa(int(consensusMessage.Sender)), 1)

	in.sendMessage(consensusMessage.Sender, rpcPair)
}

/*
	handler for propose messages
*/

func (in *Instance) handleRecorderProposeMessage(consensusMessage *proto.GenericConsensus) {
	if consensusMessage.Sender == in.getDeterministicLeader1() { //todo change this in the optimizations
		in.recorderReplicatedLog[consensusMessage.Index].P = make([]*proto.GenericConsensusValue, 1)
		in.recorderReplicatedLog[consensusMessage.Index].P[0] = &proto.GenericConsensusValue{
			Id:  consensusMessage.P[0].Id,
			Fit: strconv.FormatInt(in.Hi, 10) + "." + strconv.FormatInt(consensusMessage.Sender, 10) + "." + strconv.FormatInt(in.nodeName, 10), //fit.proposer.recorder
		}
	} else {
		randPriority := rand.Intn(int(in.Hi-10)) + 2
		in.recorderReplicatedLog[consensusMessage.Index].P = make([]*proto.GenericConsensusValue, 1)
		in.recorderReplicatedLog[consensusMessage.Index].P[0] = &proto.GenericConsensusValue{
			Id:  consensusMessage.P[0].Id,
			Fit: strconv.FormatInt(int64(randPriority), 10) + "." + strconv.FormatInt(consensusMessage.Sender, 10) + "." + strconv.FormatInt(in.nodeName, 10),
		}
	}

	fit := in.recorderReplicatedLog[consensusMessage.Index].P[0].Fit
	f1 := strings.Split(fit, ".")
	if len(f1) < 3 {
		panic("error in recorder priority assignment :priority: " + fit)
	}

	// reset the E, C sets
	in.recorderReplicatedLog[consensusMessage.Index].E = make([]*proto.GenericConsensusValue, 0)
	in.recorderReplicatedLog[consensusMessage.Index].C = make([]*proto.GenericConsensusValue, 0)

	in.debug("Recorder assigned the priority to the proposal and reset E and C", 1)
}

/*
	handler for spreadE message
*/

func (in *Instance) handleRecorderSpreadEMessage(consensusMessage *proto.GenericConsensus) {
	in.recorderReplicatedLog[consensusMessage.Index].E = in.setUnionProtoValues(in.recorderReplicatedLog[consensusMessage.Index].E, consensusMessage.P)
	in.recorderReplicatedLog[consensusMessage.Index].E = in.removeEmptyValues(in.recorderReplicatedLog[consensusMessage.Index].E)
	in.debug("Recorder processed a spreadE message and added the P set to its E set", 1)
}

/*
	handler for spreadCGatherE message
*/

func (in *Instance) handleRecorderSpreadCGatherEMessage(consensusMessage *proto.GenericConsensus) {
	in.recorderReplicatedLog[consensusMessage.Index].C = in.setUnionProtoValues(in.recorderReplicatedLog[consensusMessage.Index].C, consensusMessage.P)
	in.recorderReplicatedLog[consensusMessage.Index].C = in.removeEmptyValues(in.recorderReplicatedLog[consensusMessage.Index].C)
	in.debug("Recorder processed a SpreadCGatherE message and updated the C set", 1)
}
