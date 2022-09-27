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
	client batch store contains a map that assigns each client batch a unique identifier and allows fixed time fetching
*/

type ClientBatchStore struct {
	clientBatches sync.Map
}

/*
	Add a new client batch to the store if it is not already there
*/

func (cb *ClientBatchStore) Add(batch *proto.ClientBatch) {
	_, ok := cb.clientBatches.Load(batch.Id)
	if !ok {
		cb.clientBatches.Store(batch.Id, batch)
	}
}

/*
	return an existing client batch
*/

func (cb *ClientBatchStore) Get(id string) (*proto.ClientBatch, bool) {
	i, ok := cb.clientBatches.Load(id)
	if ok {
		return i.(*proto.ClientBatch), ok
	} else {
		return nil, ok
	}
}

/*
	Remove an element from the map
*/

func (cb *ClientBatchStore) Remove(id string) {
	cb.clientBatches.Delete(id)
}

/*
	Convert a sync.map to regular map
*/

func (cb *ClientBatchStore) convertToRegularMap(batches sync.Map) map[string]*proto.ClientBatch {
	var m map[string]*proto.ClientBatch
	m = make(map[string]*proto.ClientBatch)
	batches.Range(func(key, value interface{}) bool {
		m[fmt.Sprint(key)] = value.(*proto.ClientBatch)
		return true
	})
	return m
}

/*
	Print all the batches
*/

func (cb *ClientBatchStore) printStore(logFilePath string, nodeName int64) {
	f, err := os.Create(logFilePath + strconv.Itoa(int(nodeName)) + "mempool.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	messageBlocks := cb.convertToRegularMap(cb.clientBatches)

	for _, block := range messageBlocks {
		for j := 0; j < len(block.Messages); j++ {
			_, _ = f.WriteString(block.Messages[j].Id + ":" + block.Messages[j].Message + "\n")
		}

	}
}
