package internal

import "sync"

type MessageStore struct {
	messageBlocks map[string]messageBlock
	mutex         sync.Mutex
}

func (ms *MessageStore) init() {
	ms.messageBlocks = make(map[string]messageBlock)
}

func (ms *MessageStore) add(block messageBlock) {
	ms.mutex.Lock()
	ms.messageBlocks[block.id] = block
	ms.mutex.Unlock()
}

func (ms *MessageStore) get(id string) (messageBlock, bool) {
	ms.mutex.Lock()
	i, ok := ms.messageBlocks[id]
	ms.mutex.Unlock()
	return i, ok

}

func (ms *MessageStore) remove(id string) {
	ms.mutex.Lock()
	delete(ms.messageBlocks, id)
	ms.mutex.Unlock()
}
