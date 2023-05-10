package benchmark

import (
	"context"
	"github.com/go-redis/redis/v8"
	"raxos/proto/client"
)

/*
	struct defining the benchmark
*/

type Benchmark struct {
	mode        int           // 0 for resident k/v store, 1 for redis
	RedisClient *redis.Client // redis client
	RedisCtx    context.Context
	KVStore     map[string]string
	name        int32 // name of the server
	keyLen      int
	valueLen    int
}

/*
	initialize a new benchmark
*/

func Init(mode int, name int32, keyLen int, valueLen int) *Benchmark {

	if mode == 1 {

		client := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		rdsContext := context.Background()
		client.FlushAll(rdsContext) // delete the data base

		b := Benchmark{
			mode:        mode,
			RedisClient: client,
			RedisCtx:    rdsContext,
			KVStore:     nil,
			name:        name,
			keyLen:      keyLen,
			valueLen:    valueLen,
		}

		return &b
	} else if mode == 0 {

		b := Benchmark{
			mode:        mode,
			RedisClient: nil,
			RedisCtx:    nil,
			KVStore:     make(map[string]string),
			name:        name,
			keyLen:      keyLen,
			valueLen:    valueLen,
		}

		return &b
	} else {
		panic("should not happen")
	}

}

/*
	external API to call
*/

func (b *Benchmark) Execute(requests []client.ClientBatch) []client.ClientBatch {
	var commands []client.ClientBatch
	if b.mode == 0 {
		commands = b.residentExecute(requests)
	} else {
		commands = b.redisExecute(requests)
	}
	return commands
}

/*
	resident key value store operation: for each client request invoke the resident k/v store
*/

func (b *Benchmark) residentExecute(commands []client.ClientBatch) []client.ClientBatch {
	returnCommands := make([]client.ClientBatch, len(commands))

	for clientBatchIndex := 0; clientBatchIndex < len(commands); clientBatchIndex++ {

		returnCommands[clientBatchIndex] = client.ClientBatch{
			Id:       commands[clientBatchIndex].Id,
			Messages: make([]*client.ClientBatch_SingleMessage, len(commands[clientBatchIndex].Messages)),
			Sender:   commands[clientBatchIndex].Sender,
		}

		for clientRequestIndex := 0; clientRequestIndex < len(commands[clientBatchIndex].Messages); clientRequestIndex++ {
			returnCommands[clientBatchIndex].Messages[clientRequestIndex] = &client.ClientBatch_SingleMessage{
				Message: "",
			}

			cmd := commands[clientBatchIndex].Messages[clientRequestIndex].Message
			typ := cmd[0:1]
			key := cmd[1 : 1+b.keyLen]
			val := cmd[1+b.keyLen:]
			if typ == "0" { // write
				b.KVStore[key] = val
				returnCommands[clientBatchIndex].Messages[clientRequestIndex].Message = "0" + key + "ok"
			} else { // read
				v, ok := b.KVStore[key]
				if ok {
					returnCommands[clientBatchIndex].Messages[clientRequestIndex].Message = "1" + key + v
				} else {
					returnCommands[clientBatchIndex].Messages[clientRequestIndex].Message = "1" + key + "nil"
				}
			}
		}
	}
	return returnCommands
}

/*
	redis commands execution: batch the requests and execute
*/

func (b *Benchmark) redisExecute(commands []client.ClientBatch) []client.ClientBatch {
	returnCommands := make([]client.ClientBatch, len(commands))

	mset := make([]string, 0) // pending MSET requests
	mget := make([]string, 0) // pending MGET requests

	for clientBatchIndex := 0; clientBatchIndex < len(commands); clientBatchIndex++ {

		returnCommands[clientBatchIndex] = client.ClientBatch{
			Id:       commands[clientBatchIndex].Id,
			Messages: make([]*client.ClientBatch_SingleMessage, len(commands[clientBatchIndex].Messages)),
			Sender:   commands[clientBatchIndex].Sender,
		}

		for clientRequestIndex := 0; clientRequestIndex < len(commands[clientBatchIndex].Messages); clientRequestIndex++ {
			returnCommands[clientBatchIndex].Messages[clientRequestIndex] = &client.ClientBatch_SingleMessage{
				Message: "",
			}

			cmd := commands[clientBatchIndex].Messages[clientRequestIndex].Message
			typ := cmd[0:1]
			key := cmd[1 : 1+b.keyLen]
			val := cmd[1+b.keyLen:]
			if typ == "0" { // write
				mset = append(mset, key)
				mset = append(mset, val)
				returnCommands[clientBatchIndex].Messages[clientRequestIndex].Message = "0" + key + "ok" // writes always succeed
			} else { // read
				mget = append(mget, key)
				returnCommands[clientBatchIndex].Messages[clientRequestIndex].Message = ""
			}
		}
	}

	if len(mset) > 0 {
		// execute writes in a batch
		if err := b.RedisClient.MSet(b.RedisCtx, mset).Err(); err != nil {
			panic(err)
		}
	}

	if len(mget) > 0 {
		// execute reads in a batch
		vs, err := b.RedisClient.MGet(b.RedisCtx, mget...).Result()
		if err != nil {
			panic(err)
		}

		vsCount := 0

		for clientBatchIndex := 0; clientBatchIndex < len(commands); clientBatchIndex++ {
			for clientRequestIndex := 0; clientRequestIndex < len(commands[clientBatchIndex].Messages); clientRequestIndex++ {
				cmd := commands[clientBatchIndex].Messages[clientRequestIndex].Message
				typ := cmd[0:1]
				key := cmd[1 : 1+b.keyLen]
				if typ == "0" {
					// we already set the response for writes
				} else { // read
					if vs[vsCount] == nil {
						// key not found
						returnCommands[clientBatchIndex].Messages[clientRequestIndex].Message = "1" + key + "nil"
					} else {
						if rep, ok := vs[vsCount].(string); !ok {
							panic(vs[vsCount])
						} else {
							returnCommands[clientBatchIndex].Messages[clientRequestIndex].Message = "1" + key + rep
						}
					}
					vsCount++
				}
			}
		}
	}

	return returnCommands
}
