package benchmark

import (
	"strings"
)

type KVStoreApplication struct {
	register map[string]string
	numKeys  int
}

func (kvStoreApp *KVStoreApplication) executeKVStoreApp(request string) {
	// request pattern: UPDATE/READ:key:value
	op := strings.Split(request, ":")[0]
	key := strings.Split(request, ":")[1]
	value := strings.Split(request, ":")[2]
	if op == "READ" {
		kvStoreApp.get(key)
	} else if op == "UPDATE" {
		kvStoreApp.put(key, value)
	}
}

func (kvStoreApp *KVStoreApplication) Init(numKeys int) {
	kvStoreApp.numKeys = numKeys
	kvStoreApp.register = make(map[string]string)
	for i := 0; i < numKeys; i++ {
		kvStoreApp.register[Get23LengthRecord(i)] = Get100LengthValue()
	}
}

func (kvStoreApp *KVStoreApplication) get(key string) string {
	return kvStoreApp.register[key]
}

func (kvStoreApp *KVStoreApplication) put(key string, value string) {
	kvStoreApp.register[key] = value
}
