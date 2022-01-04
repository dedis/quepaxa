package benchmark

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strings"
	"time"
)

const redisBatchSize = 50
const redisBatchTime = 100 // microseconds

type RedisApplication struct {
	numKeys int
	client  *redis.Client
	ctx     context.Context
	ops     chan string
}

func (redisApp *RedisApplication) Init(numKeys int) {
	redisApp.numKeys = numKeys
	redisApp.ctx = context.Background()
	redisApp.client = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	redisApp.ops = make(chan string, 1000)
	redisApp.client.FlushAll(redisApp.ctx)

	for i := 0; i < redisApp.numKeys; i++ {
		redisApp.client.HSet(redisApp.ctx, Get23LengthRecord(i),
			"field0", Get100LengthValue(),
			"field1", Get100LengthValue(),
			"field2", Get100LengthValue(),
			"field3", Get100LengthValue(),
			"field4", Get100LengthValue(),
			"field5", Get100LengthValue(),
			"field6", Get100LengthValue(),
			"field7", Get100LengthValue(),
			"field8", Get100LengthValue(),
			"field9", Get100LengthValue())
	}

	go redisApp.run()

}

func (redisApp *RedisApplication) executeRedisApp(request string) {
	redisApp.ops <- request
}

//func (redisApp *RedisApplication) get(key string) map[string]string {
//	rslt, _ := redisApp.client.HGetAll(redisApp.ctx, key).Result()
//	return rslt
//}
//
//func (redisApp *RedisApplication) update(key string, field string, value string) {
//	redisApp.client.HSet(redisApp.ctx, key, field, value)
//}

func (redisApp *RedisApplication) run() {
	lastCommitted := time.Now()
	for true {
		numRequests := int64(0)
		var requests []string
		for !(numRequests >= redisBatchSize || (time.Now().Sub(lastCommitted).Microseconds() > redisBatchTime && numRequests > 0)) {
			request := <-redisApp.ops
			requests = append(requests, request)
			numRequests++
		}
		pipe := redisApp.client.TxPipeline()
		for i := 0; i < len(requests); i++ {
			request := requests[i]
			// request pattern: UPDATE/READ:key:field:value
			op := strings.Split(request, ":")[0]
			key := strings.Split(request, ":")[1]
			field := strings.Split(request, ":")[2]
			value := strings.Split(request, ":")[3]
			if op == "READ" {
				pipe.HGetAll(redisApp.ctx, key)
			} else if op == "UPDATE" {
				pipe.HSet(redisApp.ctx, key, field, value)
			}
		}
		_, _ = pipe.Exec(redisApp.ctx)
		lastCommitted = time.Now()
	}
}
