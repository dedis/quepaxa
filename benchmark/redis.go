package benchmark

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strings"
	"time"
)

/*
todo this implementation does not support replies. To support replies it has to invoke the driver for each request, and it doesn't give goo performance
todo Check the Rabia Redis integration
*/

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
		redisApp.client.HSet(redisApp.ctx, GetNLengthRecord(i, 23),
			"field0", GetNLengthValue(100),
			"field1", GetNLengthValue(100),
			"field2", GetNLengthValue(100),
			"field3", GetNLengthValue(100),
			"field4", GetNLengthValue(100),
			"field5", GetNLengthValue(100),
			"field6", GetNLengthValue(100),
			"field7", GetNLengthValue(100),
			"field8", GetNLengthValue(100),
			"field9", GetNLengthValue(100))
	}

	go redisApp.run()

}

func (redisApp *RedisApplication) executeRedisApp(request string) string {
	redisApp.ops <- request
	return "success"
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
