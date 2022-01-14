package raxos

import (
	"fmt"
	"raxos/proto"
	"sync"
)

/*

Message store contains a map that assigns each block (batch of batches of client requests) a unique identifier and allows constant time fetching

*/

type Block struct {
	messageBlock *proto.MessageBlock // to avoid unnecessary conversions we use the same proto Message type for storing a block
	acks         []int64             // contains the set of nodes who acknowledged the block
}

/*

Message store object will be accessed by the main thread and the updateStateMachine thread. Since go maps are not thread safe, we use a mutex
*/

type MessageStore struct {
	messageBlocks map[string]Block
	mutex         sync.Mutex
}

func (ms *MessageStore) Init() {
	ms.messageBlocks = make(map[string]Block)
}

/*

Adds a new block to the store if its not already there
*/

func (ms *MessageStore) Add(block *proto.MessageBlock) {
	ms.mutex.Lock()
	_, ok := ms.messageBlocks[block.Hash]
	if !ok {
		ms.messageBlocks[block.Hash] = Block{
			messageBlock: block,
		}
	}
	ms.mutex.Unlock()
}

/*return an existing block*/

func (ms *MessageStore) Get(id string) (*proto.MessageBlock, bool) {
	ms.mutex.Lock()
	i, ok := ms.messageBlocks[id]
	ms.mutex.Unlock()
	return i.messageBlock, ok
}

/*return the set of acks for a given block*/

func (ms *MessageStore) getAcks(id string) []int64 {
	ms.mutex.Lock()
	_, ok := ms.messageBlocks[id]
	ms.mutex.Unlock()
	if ok {
		return ms.messageBlocks[id].acks
	}
	return nil

}

/* Remove an element from the map*/

func (ms *MessageStore) Remove(id string) {
	ms.mutex.Lock()
	delete(ms.messageBlocks, id)
	ms.mutex.Unlock()
}

/*add a new ack to the ack list of a block*/

func (ms *MessageStore) addAck(id string) {
	ms.mutex.Lock()
	_, ok := ms.messageBlocks[id]
	ms.mutex.Unlock()
	if ok {
		ms.mutex.Lock()
		tempAcks := ms.messageBlocks[id].acks
		tempBlock := ms.messageBlocks[id].messageBlock
		tempAcks = append(tempAcks, 1)
		ms.messageBlocks[id] = Block{
			messageBlock: tempBlock,
			acks:         tempAcks,
		}
		ms.mutex.Unlock()
	}
}

/*Print all the blocks*/

func (ms *MessageStore) printStore() {
	ms.mutex.Lock()
	for hash, block := range ms.messageBlocks {
		fmt.Print(hash, ":")
		for i := 0; i < len(block.messageBlock.Requests); i++ {
			for j := 0; j < len(block.messageBlock.Requests[i].Requests); j++ {
				fmt.Print(block.messageBlock.Requests[i].Requests[j].Message, ",")
			}
		}
		fmt.Print("\n")
	}
	ms.mutex.Unlock()
}
