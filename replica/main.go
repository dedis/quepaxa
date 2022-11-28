package main

import (
	"flag"
	"fmt"
	"os"
	"raxos/configuration"
	raxos "raxos/replica/src"
	"time"
)

func main() {
	configFile := flag.String("config", "configuration/local/configuration.yml", "raxos configuration file")
	name := flag.Int64("name", 1, "name of the replica")
	logFilePath := flag.String("logFilePath", "logs/", "log file path")
	batchSize := flag.Int64("batchSize", 50, "replica batch size")
	batchTime := flag.Int64("batchTime", 2, "replica batch time in milli seconds")
	leaderTimeout := flag.Int64("leaderTimeout", 500, "leader timeout in milli seconds")
	pipelineLength := flag.Int64("pipelineLength", 50, "pipeline length maximum number of outstanding proposals")
	benchmark := flag.Int64("benchmark", 0, "Benchmark: 0 for echo service, 1 for KV store and 2 for Redis ")
	debugOn := flag.Bool("debugOn", false, "true / false")
	debugLevel := flag.Int("debugLevel", 0, "debug level")
	leaderMode := flag.Int("leaderMode", 0, "mode of leader change: 0 for fixed leader order")
	serverMode := flag.Int("serverMode", 0, "0 for non-lan-optimized, 1 for lan optimized")
	flag.Parse()

	cfg, err := configuration.NewInstanceConfig(*configFile, *name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}

	serverInstance := raxos.New(cfg, *name, *logFilePath, *batchSize, *leaderTimeout, *pipelineLength, *benchmark, *debugOn, *debugLevel, *leaderMode, *serverMode, *batchTime) // create a new server instance

	serverInstance.NetworkInit()
	serverInstance.Run()

	/*to avoid exiting the main thread*/
	for true {
		time.Sleep(10 * time.Second)
	}
}
