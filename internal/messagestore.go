package internal

type MessageStore struct {
	messageBlocks map[string]MessageBlock
	//mutex         sync.Mutex
}

func (ms *MessageStore) Init() {
	ms.messageBlocks = make(map[string]MessageBlock)
}

func (ms *MessageStore) Add(block MessageBlock) {
	//ms.mutex.Lock()
	ms.messageBlocks[block.id] = block
	//ms.mutex.Unlock()
}

func (ms *MessageStore) Get(id string) (MessageBlock, bool) {
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
