package benchmark

import (
	"math/rand"
	"strconv"
)

type App struct {
	Workload   int64
	NoOpApp    *NoOpApplication
	RedisApp   *RedisApplication
	KvstoreApp *KVStoreApplication
}

func (app *App) Process(request string) {
	if app.Workload == 0 {
		app.NoOpApp.executeNoOpApp(request)
	}
	if app.Workload == 1 {
		app.KvstoreApp.executeKVStoreApp(request)
	}
	if app.Workload == 2 {
		app.RedisApp.executeRedisApp(request)
	}
}

func Get100LengthValue() string {
	str := strconv.Itoa(rand.Intn(10))
	size := len(str)
	for size <= 100 {
		str = strconv.Itoa(rand.Intn(10)) + str
		size = len(str)
	}
	return str
}

func Get23LengthRecord(i int) string {
	str := strconv.Itoa(i)
	size := len(str)
	for size <= 19 {
		str = "0" + str
		size = len(str)
	}
	return "user" + str
}

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
