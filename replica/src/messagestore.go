package raxos

import (
	"fmt"
	"log"
	"os"
	"raxos/proto"
	"strconv"
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
	Message store object will be accessed by the main thread and the updateStateMachine thread. Since go maps are not thread safe, we use a sync.maps
*/

type MessageStore struct {
	messageBlocks sync.Map
}

/*
	Sync.map doesn't need to be allocated
*/

func (ms *MessageStore) Init() {
}

/*
	Add a new block to the store if it is not already there
*/

func (ms *MessageStore) Add(block *proto.MessageBlock) {
	_, ok := ms.messageBlocks.Load(block.Hash)
	if !ok {
		ms.messageBlocks.Store(block.Hash, Block{
			messageBlock: block,
		})
	}
}

/*
	return an existing block
*/

func (ms *MessageStore) Get(id string) (*proto.MessageBlock, bool) {
	i, ok := ms.messageBlocks.Load(id)
	if ok {
		return i.(Block).messageBlock, ok
	} else {
		return nil, ok
	}
}

/*
	return the set of acks for a given block
*/

func (ms *MessageStore) getAcks(id string) []int64 {
	block, ok := ms.messageBlocks.Load(id)
	if ok {
		return block.(Block).acks
	}
	return nil
}

/*
	Remove an element from the map
*/

func (ms *MessageStore) Remove(id string) {
	ms.messageBlocks.Delete(id)
}

/*
	add a new ack to the ack list of a block
*/

func (ms *MessageStore) addAck(id string) {
	block, ok := ms.messageBlocks.Load(id)
	if ok {
		tempAcks := block.(Block).acks
		tempBlock := block.(Block).messageBlock
		tempAcks = append(tempAcks, 1)
		ms.messageBlocks.Delete(id)
		ms.messageBlocks.Store(id, Block{
			messageBlock: tempBlock,
			acks:         tempAcks,
		})
	}
}

/*
	Print all the blocks
*/

func (ms *MessageStore) printStore(logFilePath string, nodeName int64) {
	f, err := os.Create(logFilePath + strconv.Itoa(int(nodeName)) + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	messageBlocks := ms.convertToRegularMap(ms.messageBlocks)

	for hash, block := range messageBlocks {
		for i := 0; i < len(block.messageBlock.Requests); i++ {
			_, _ = f.WriteString(hash + "-" + block.messageBlock.Requests[i].Id + ":")
			for j := 0; j < len(block.messageBlock.Requests[i].Requests); j++ {
				_, _ = f.WriteString(block.messageBlock.Requests[i].Requests[j].Message + ",")
			}
			_, _ = f.WriteString("\n")
		}
	}
}

/*

	Convert a sync.map to regular map
*/

func (ms *MessageStore) convertToRegularMap(blocks sync.Map) map[string]Block {
	var m map[string]Block
	m = make(map[string]Block)
	blocks.Range(func(key, value interface{}) bool {
		m[fmt.Sprint(key)] = value.(Block)
		return true
	})
	return m
}
