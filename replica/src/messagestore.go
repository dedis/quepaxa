package raxos

import (
	"fmt"
	"raxos/proto"
)

type Block struct {
	messageBlock *proto.MessageBlock
	acks         []int64
}

type MessageStore struct {
	messageBlocks map[string]Block
	//mutex         sync.Mutex
}

func (ms *MessageStore) Init() {
	ms.messageBlocks = make(map[string]Block)
}

func (ms *MessageStore) Add(block *proto.MessageBlock) {
	//ms.mutex.Lock()
	_, ok := ms.messageBlocks[block.Hash]
	if !ok {
		ms.messageBlocks[block.Hash] = Block{
			messageBlock: block,
		}
	}
	//ms.mutex.Unlock()
}

func (ms *MessageStore) Get(id string) (*proto.MessageBlock, bool) {
	//ms.mutex.Lock()
	i, ok := ms.messageBlocks[id]
	//ms.mutex.Unlock()
	return i.messageBlock, ok
}

func (ms *MessageStore) getAcks(id string) []int64 {
	//ms.mutex.Lock()
	_, ok := ms.messageBlocks[id]
	if ok {
		return ms.messageBlocks[id].acks
	}
	return nil
	//ms.mutex.Unlock()
}

func (ms *MessageStore) Remove(id string) {
	//ms.mutex.Lock()
	delete(ms.messageBlocks, id)
	//ms.mutex.Unlock()
}

func (ms *MessageStore) addAck(id string) {
	_, ok := ms.messageBlocks[id]
	if ok {
		tempAcks := ms.messageBlocks[id].acks
		tempBlock := ms.messageBlocks[id].messageBlock
		tempAcks = append(tempAcks, 1)
		ms.messageBlocks[id] = Block{
			messageBlock: tempBlock,
			acks:         tempAcks,
		}
	}
}

func (ms *MessageStore) printStore() {
	for hash, block := range ms.messageBlocks {
		fmt.Print("Hash:", hash, ": \n")
		for i := 0; i < len(block.messageBlock.Requests); i++ {
			for j := 0; j < len(block.messageBlock.Requests[i].Requests); j++ {
				fmt.Print(block.messageBlock.Requests[i].Requests[j].Message, " \n")
			}
		}
		fmt.Print("\n")
	}
}
