package benchmark

import (
	"math/rand"
	"strconv"
)

/*

	App defines a generic state machine. currently, it supports three implementations

	(1) no-op app: an echo app which returns the request as the response, with an added delay
	(2) key value store app: a key value store
	(3) redis key value store
*/

type App struct {
	Workload   int64
	NoOpApp    *NoOpApplication
	RedisApp   *RedisApplication
	KvstoreApp *KVStoreApplication
}

/*
	Process is the generic interface for each type of App. It calls the registered app for each request
*/

func (app *App) Process(request string) string {
	if app.Workload == 0 {
		return app.NoOpApp.executeNoOpApp(request)
	}
	if app.Workload == 1 {
		return app.KvstoreApp.executeKVStoreApp(request)
	}
	if app.Workload == 2 {
		return app.RedisApp.executeRedisApp(request)
	}
	return "error"
}

/*
	Generate a N random length string which consists of 1-9
*/

func GetNLengthValue(N int) string {
	str := strconv.Itoa(rand.Intn(10))
	size := len(str)
	for size <= N {
		str = strconv.Itoa(rand.Intn(10)) + str
		size = len(str)
	}
	return str
}

/*
	Generate a N length record of the form user..0..n
*/

func GetNLengthRecord(i int, N int) string {
	str := strconv.Itoa(i)
	size := len(str)
	for size <= N-4 {
		str = "0" + str
		size = len(str)
	}
	return "user" + str
}

/*
	Initialize a new app instance, and init() app if required
*/

func InitApp(b int64, serviceTime int64, numKeys int64) *App {
	noOpApp := NoOpApplication{SleepDuration: serviceTime}
	kvStoreApp := KVStoreApplication{}
	redisApp := RedisApplication{}

	if b == 1 {
		kvStoreApp.Init(int(numKeys))
	}
	if b == 2 {
		redisApp.Init(int(numKeys))
	}
	app := App{
		Workload:   b,
		NoOpApp:    &noOpApp,
		RedisApp:   &redisApp,
		KvstoreApp: &kvStoreApp,
	}

	return &app
}
