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
	serviceTime := flag.Int64("serviceTime", 1, "service time in micro seconds")
	responseSize := flag.Int64("responseSize", 8, "response size in bytes")
	batchSize := flag.Int64("batchSize", 50, "replica batch size")
	batchTime := flag.Int64("batchTime", 50, "maximum time to wait for collecting a batch of requests in micro seconds")
	leaderTimeout := flag.Int64("leaderTimeout", 50, "leader timeout in milli seconds")
	pipelineLength := flag.Int64("pipelineLength", 50, "pipeline length maximum number of outstanding proposals")
	benchmark := flag.Int64("benchmark", 0, "Benchmark: 0 for echo service, 1 for KV store and 2 for Redis ")
	numKeys := flag.Int64("numKeys", 1000, "Number of keys in the key value store")
	hashBatchSize := flag.Int64("hashBatchSize", 10, "replica hash batch size")
	hashBatchTime := flag.Int64("hashBatchTime", 1000, "maximum time to wait for collecting a batch of hashes in micro seconds")
	debugOn := flag.Bool("debugOn", false, "true / false")
	debugLevel := flag.Int("debugLevel", 0, "debug level")
	flag.Parse()

	cfg, err := configuration.NewInstanceConfig(*configFile, *name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}

	in := raxos.New(cfg, *name, *logFilePath, *serviceTime, *responseSize, *batchSize, *batchTime, *leaderTimeout, *pipelineLength, *benchmark, *numKeys, *hashBatchSize, *hashBatchTime, *debugOn, *debugLevel)

	in.Run()
	in.StartOutgoingLinks()
	in.CollectAndProposeHashes()
	in.BroadcastBlock()
	go in.WaitForConnections()

	/*to avoid exiting the main thread*/
	for true {
		time.Sleep(10 * time.Second)
	}
}
