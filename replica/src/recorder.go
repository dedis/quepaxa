package raxos

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
	"raxos/configuration"
	"raxos/proto/client"
	"strconv"
	"sync"
	"time"
)

type Value struct {
	priority    int64
	proposer_id int64
	thread_id   int64
	ids         []string
}

// Recorder side slot

type RecorderSlot struct {
	Mutex *sync.Mutex
	S     int
	F     Value
	A     Value
	M     Value
}

type Recorder struct {
	address               string       // address to listen for gRPC connections
	listener              net.Listener // socket for gRPC connections
	server                *grpc.Server // gRPC server
	connection            *GRPCConnection
	clientBatches         *ClientBatchStore
	lastSeenTimeProposers []*time.Time // last seen times of each proposer
	recorderToProxyChan   chan Decision
	name                  int64
	slots                 []RecorderSlot // recorder side replicated log
	cfg                   configuration.InstanceConfig
	instanceCreationMutex *sync.Mutex
	debugOn               bool // if turned on, the debug messages will be print on the console
	debugLevel            int  // debug level
}

// instantiate a new Recorder

func NewRecorder(cfg configuration.InstanceConfig, clientBatches *ClientBatchStore, lastSeenTimeProposers []*time.Time, recorderToProxyChan chan Decision, name int64, debugOn bool, debugLevel int) *Recorder {

	re := Recorder{
		address:               "",
		listener:              nil,
		server:                nil,
		connection:            nil,
		clientBatches:         clientBatches,
		lastSeenTimeProposers: lastSeenTimeProposers,
		recorderToProxyChan:   recorderToProxyChan,
		name:                  name,
		slots:                 make([]RecorderSlot, 0),
		cfg:                   cfg,
		instanceCreationMutex: &sync.Mutex{},
		debugLevel:            debugLevel,
		debugOn:               debugOn,
	}

	// serverAddress
	for i := 0; i < len(cfg.Peers); i++ {
		intName, _ := strconv.Atoi(cfg.Peers[i].Name)
		if re.name == int64(intName) {
			re.address = "0.0.0.0:" + cfg.Peers[i].RECORDERPORT
			break
		}
	}

	return &re
}

/*
	if turned on, print the message to console
*/

func (re *Recorder) debug(message string, level int) {
	if re.debugOn && level >= re.debugLevel {
		fmt.Printf("%s\n", message)
	}
}

// start listening to gRPC connection

func (r *Recorder) NetworkInit() {
	r.server = grpc.NewServer()
	r.connection = &GRPCConnection{Recorder: r}
	RegisterConsensusServer(r.server, r.connection)

	// start listener
	listener, err := net.Listen("tcp", r.address)
	if err != nil {
		panic("listen: %v")
	}
	r.listener = listener
	go func() {
		err := r.server.Serve(listener)
		if err != nil {
			panic("should not happen")
		}
	}()
}

// check if all the batches are available in the store

func (re *Recorder) findAllBatches(ids []string) bool {
	for i := 0; i < len(ids); i++ {
		_, ok := re.clientBatches.Get(ids[i])
		if !ok {
			return false
		}
	}
	return true
}

// return the max of oldValue and the new value

func (r *Recorder) max(oldValue Value, p *ProposerMessage_Proposal) Value {
	maxm := oldValue
	if maxm.priority < p.Priority {
		maxm = Value{
			priority:    p.Priority,
			proposer_id: p.ProposerId,
			thread_id:   p.ThreadId,
			ids:         p.Ids,
		}
		return maxm
	} else if maxm.priority == p.Priority {
		if maxm.proposer_id < p.ProposerId {
			maxm = Value{
				priority:    p.Priority,
				proposer_id: p.ProposerId,
				thread_id:   p.ThreadId,
				ids:         p.Ids,
			}
			return maxm
		} else if maxm.proposer_id == p.ProposerId {
			if maxm.thread_id < p.ThreadId {
				maxm = Value{
					priority:    p.Priority,
					proposer_id: p.ProposerId,
					thread_id:   p.ThreadId,
					ids:         p.Ids,
				}
				return maxm
			} else if maxm.thread_id == p.ThreadId {
				panic("should not happen")
			}
		}
	}

	return maxm
}

// main recorder logic goes here

func (re *Recorder) espImpl(index int64, s int, p *ProposerMessage_Proposal) (int64, Value, Value) {

	re.instanceCreationMutex.Lock()
	for int64(len(re.slots)) < index+1 {
		re.slots = append(re.slots, RecorderSlot{
			Mutex: &sync.Mutex{},
			S:     0,
			F: Value{
				priority:    -1,
				proposer_id: -1,
				thread_id:   -1,
				ids:         nil,
			},
			A: Value{
				priority:    -1,
				proposer_id: -1,
				thread_id:   -1,
				ids:         nil,
			},
			M: Value{
				priority:    -1,
				proposer_id: -1,
				thread_id:   -1,
				ids:         nil,
			},
		})
	}
	re.instanceCreationMutex.Unlock()

	re.slots[index].Mutex.Lock()
	if re.slots[index].S == s {
		re.slots[index].A = re.max(re.slots[index].A, p)
	} else if re.slots[index].S < s {
		if re.slots[index].S+1 < s {
			re.slots[index].A = Value{
				priority:    -1,
				proposer_id: -1,
				thread_id:   -1,
				ids:         nil,
			}
		}
		re.slots[index].S = s
		re.slots[index].F = Value{
			priority:    p.Priority,
			proposer_id: p.ProposerId,
			thread_id:   p.ThreadId,
			ids:         p.Ids,
		}
		re.slots[index].M = re.slots[index].A
		re.slots[index].A = Value{
			priority:    p.Priority,
			proposer_id: p.ProposerId,
			thread_id:   p.ThreadId,
			ids:         p.Ids,
		}
	}

	returnS := re.slots[index].S
	returnF := re.slots[index].F
	returnM := re.slots[index].M

	re.slots[index].Mutex.Unlock()

	return int64(returnS), returnF, returnM

}

// util method to convert between two prototypes

func (r *Recorder) convertToClientBatchMessages(messages []*ProposerMessage_ClientBatch_SingleMessage) []*client.ClientBatch_SingleMessage {
	ra := make([]*client.ClientBatch_SingleMessage, len(messages))
	for i := 0; i < len(messages); i++ {
		ra[i] = &client.ClientBatch_SingleMessage{
			Message: messages[i].Message,
		}
	}
	return ra
}

// answer to proposer RPC

func (re *Recorder) HandleESP(req *ProposerMessage) *RecorderResponse {

	var response RecorderResponse

	// send the last decided index details to the proxy, if available
	if len(req.DecidedSlots) > 0 {
		d := Decision{
			indexes:   make([]int, 0),
			decisions: make([][]string, 0),
		}

		for i := 0; i < len(req.DecidedSlots); i++ {
			d.indexes = append(d.indexes, int(req.DecidedSlots[i].Index))
			d.decisions = append(d.decisions, req.DecidedSlots[i].Ids)
		}

		re.recorderToProxyChan <- d
	}

	if len(req.P.ClientBatches) == 0 {
		// if there are only hashes, then check if all the client batches are available in the shared pool
		allBatchesFound := re.findAllBatches(req.P.Ids)
		if !allBatchesFound {
			response.HasClientBacthes = false
			return &response
		}
	}

	if len(req.P.ClientBatches) > 0 {
		// add all the batches to the store
		for i := 0; i < len(req.P.ClientBatches); i++ {
			re.clientBatches.Add(client.ClientBatch{
				Sender:   req.P.ClientBatches[i].Sender,
				Messages: re.convertToClientBatchMessages(req.P.ClientBatches[i].Messages),
				Id:       req.P.ClientBatches[i].Id,
			})
		}
	}

	// process using the recorder logic
	S, F, M := re.espImpl(req.Index, int(req.S), req.P)
	response.S = S
	response.F = &RecorderResponse_Proposal{
		Priority:   F.priority,
		ProposerId: F.proposer_id,
		ThreadId:   F.thread_id,
		Ids:        F.ids,
	}
	response.M = &RecorderResponse_Proposal{
		Priority:   M.priority,
		ProposerId: M.proposer_id,
		ThreadId:   M.thread_id,
		Ids:        M.ids,
	}

	// Mark the time of the proposal message for the proposer
	proposer := req.Sender
	*re.lastSeenTimeProposers[proposer] = time.Now()

	return &response
}

// convert between proto types

func (r *Recorder) convertToDecideResponseClientBatchMessages(messages []*client.ClientBatch_SingleMessage) []*DecideResponse_ClientBatch_SingleMessage {
	rt := make([]*DecideResponse_ClientBatch_SingleMessage, 0)
	for i := 0; i < len(messages); i++ {
		rt = append(rt, &DecideResponse_ClientBatch_SingleMessage{
			Message: messages[i].Message,
		})
	}

	return rt
}

// answer to fetch request

func (r *Recorder) HandleFetch(req *DecideRequest) *DecideResponse {
	response := DecideResponse{
		ClientBatches: nil,
	}
	response.ClientBatches = make([]*DecideResponse_ClientBatch, 0)
	for i := 0; i < len(req.Ids); i++ {
		btch, ok := r.clientBatches.Get(req.Ids[i])
		if ok {
			response.ClientBatches = append(response.ClientBatches, &DecideResponse_ClientBatch{
				Sender:   btch.Sender,
				Messages: r.convertToDecideResponseClientBatchMessages(btch.Messages),
				Id:       btch.Id,
			})
		}
	}

	return &response
}
