package main

import (
	"flag"
	"fmt"
	"os"
	"raxos/configuration"
	raxos "raxos/replica/src"
	"time"
)

// this file defines the main routine of QuePaxa, which takes input arguments from the command line

func main() {
	configFile := flag.String("config", "configuration/local/configuration.yml", "QuePaxa configuration file")
	name := flag.Int64("name", 1, "name of the replica")
	logFilePath := flag.String("logFilePath", "logs/", "log file path")
	batchSize := flag.Int64("batchSize", 50, "replica batch size")
	batchTime := flag.Int64("batchTime", 5000, "replica batch time in micro seconds")
	leaderTimeout := flag.Int64("leaderTimeout", 5000000, "leader timeout in micro seconds")
	pipelineLength := flag.Int64("pipelineLength", 1, "pipeline length maximum number of outstanding proposals")
	benchmark := flag.Int64("benchmark", 0, "Benchmark: 0 for KV store and 1 for Redis ")
	debugOn := flag.Bool("debugOn", false, "true / false")
	isAsync := flag.Bool("isAsync", false, "true / false to simulate consensus level asynchrony")
	debugLevel := flag.Int("debugLevel", 1010, "debug level")
	leaderMode := flag.Int("leaderMode", 0, "mode of leader change: 0 for fixed leader order, 1 for round robin, static partition,  2 for M.A.B based on commit times, 3 for asynchronous, 4 for last committed proposer")
	serverMode := flag.Int("serverMode", 0, "0 for non-lan-optimized, 1 for lan optimized")
	epochSize := flag.Int("epochSize", 100, "epoch size for MAB")
	keyLen := flag.Int64("keyLen", 8, "length of key")
	valLen := flag.Int64("valLen", 8, "length of value")
	asyncTimeOut := flag.Int64("asyncTimeOut", 500, "artificial async timeout in milli seconds")
	requestPropagationTime := flag.Int64("requestPropagationTime", 0, "additional wait time in 'milli seconds' for client batches, such that there is enough time for client driven request propagation, for LAN this is 0, for WAN it is usually set to 5ms")
	flag.Parse()

	cfg, err := configuration.NewInstanceConfig(*configFile, *name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}

	serverInstance := raxos.New(cfg, *name, *logFilePath, *batchSize, *leaderTimeout, *pipelineLength, *debugOn, *debugLevel, *leaderMode, *serverMode, *batchTime, *epochSize, int(*benchmark), int(*keyLen), int(*valLen), *requestPropagationTime, *isAsync, *asyncTimeOut) // create a new server instance

	serverInstance.NetworkInit()
	serverInstance.Run()

	/*to avoid exiting the main thread*/
	for true {
		time.Sleep(10 * time.Second)
	}
}
