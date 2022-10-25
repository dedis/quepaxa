package raxos

import (
	"fmt"
	"log"
	"os"
	"raxos/proto/client"
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
	add a new client batch to the store if it is not already there
*/

func (cb *ClientBatchStore) Add(batch client.ClientBatch) {
	_, ok := cb.clientBatches.Load(batch.Id)
	if !ok {
		cb.clientBatches.Store(batch.Id, batch)
	}
}

/*
	return an existing client batch
*/

func (cb *ClientBatchStore) Get(id string) (client.ClientBatch, bool) {
	i, ok := cb.clientBatches.Load(id)
	if ok {
		return i.(client.ClientBatch), ok
	} else {
		return client.ClientBatch{}, ok
	}
}

/*
	remove an element from the map
*/

func (cb *ClientBatchStore) Remove(id string) {
	cb.clientBatches.Delete(id)
}

/*
	convert a sync.map to regular map
*/

func (cb *ClientBatchStore) convertToRegularMap() map[string]client.ClientBatch {
	var m map[string]client.ClientBatch
	m = make(map[string]client.ClientBatch)
	cb.clientBatches.Range(func(key, value interface{}) bool {
		m[fmt.Sprint(key)] = value.(client.ClientBatch)
		return true
	})
	return m
}

/*
	print all the batches
*/

func (cb *ClientBatchStore) printStore(logFilePath string, nodeName int64) {
	f, err := os.Create(logFilePath + strconv.Itoa(int(nodeName)) + "-mempool.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	messageBlocks := cb.convertToRegularMap()

	for id, block := range messageBlocks {
		for j := 0; j < len(block.Messages); j++ {
			_, _ = f.WriteString(id + "." + strconv.Itoa(j) + ":" + block.Messages[j].Message + "\n")
		}

	}
}
