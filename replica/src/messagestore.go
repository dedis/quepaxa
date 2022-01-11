package raxos

import (
	"fmt"
	"raxos/proto"
)

type MessageStore struct {
	messageBlocks map[string]*proto.MessageBlock
	//mutex         sync.Mutex
}

func (ms *MessageStore) Init() {
	ms.messageBlocks = make(map[string]*proto.MessageBlock)
}

func (ms *MessageStore) Add(block *proto.MessageBlock) {
	//ms.mutex.Lock()
	ms.messageBlocks[block.Hash] = block
	//ms.mutex.Unlock()
}

func (ms *MessageStore) Get(id string) (*proto.MessageBlock, bool) {
	//ms.mutex.Lock()
	i, ok := ms.messageBlocks[id]
	//ms.mutex.Unlock()
	return i, ok

}

func (ms *MessageStore) Remove(id string) {
	//ms.mutex.Lock()
	delete(ms.messageBlocks, id)
	//ms.mutex.Unlock()
}

func (ms *MessageStore) printStore() {
	for hash, block := range ms.messageBlocks {
		fmt.Print("Hash:", hash, ": \n")
		for i := 0; i < len(block.Requests); i++ {
			for j := 0; j < len(block.Requests[i].Requests); j++ {
				fmt.Print(block.Requests[i].Requests[j].Message, " \n")
			}
		}
		fmt.Print("\n")
	}
}
