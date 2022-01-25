package benchmark

import (
	"strings"
)

/*
	KVStoreApplication implements a key value store. The key value store is initialized with a fixed set of keys (num Keys). It supports get and put operations
	To make sure that the client reads correspond to an existing key, the numKeys parameter in both the client and the replica should be the same
*/

type KVStoreApplication struct {
	register map[string]string
	numKeys  int
	/*
		No need to have a mutex because SMR is touched only by a single thread
	*/
}

/*
	Populate the key value store with initial values
	23 key size and 100 length value sizes are taken from the YCSB-A workload
	// todo might want to make the key lengths and value lengths different, depending on the corresponding values of Rabia and EPaxos
*/

func (kvStoreApp *KVStoreApplication) Init(numKeys int) {
	kvStoreApp.numKeys = numKeys
	kvStoreApp.register = make(map[string]string)
	for i := 0; i < numKeys; i++ {
		kvStoreApp.register[GetNLengthRecord(i, 23)] = GetNLengthValue(100)
	}
}

/*
	request: UPDATE/READ:key:value
	divide the request into sections and execute it
*/

func (kvStoreApp *KVStoreApplication) executeKVStoreApp(request string) string {
	// request pattern: UPDATE/READ:key:value
	op := strings.Split(request, ":")[0]
	key := strings.Split(request, ":")[1]
	value := strings.Split(request, ":")[2]
	if op == "READ" {
		return kvStoreApp.get(key)
	} else if op == "UPDATE" {
		kvStoreApp.put(key, value)
		return "success"
	}
	return ""
}

/*
	Helper functions to get a record corresponding to a key
*/

func (kvStoreApp *KVStoreApplication) get(key string) string {
	return kvStoreApp.register[key]
}

/*
	Helper function to put a new {key:value}
*/

func (kvStoreApp *KVStoreApplication) put(key string, value string) {
	kvStoreApp.register[key] = value
}
