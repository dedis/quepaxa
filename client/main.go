package main

import (
	"flag"
	"fmt"
	"os"
	"raxos/client/cmd"
	"raxos/configuration"
)

func main() {
	name := flag.Int64("name", 5, "name of the client as specified in the configuration.yml")
	configFile := flag.String("config", "configuration/local/configuration.yml", "raxos configuration file")
	logFilePath := flag.String("logFilePath", "logs/", "log file path")
	batchSize := flag.Int64("batchSize", 50, "client batch size")
	batchTime := flag.Int64("batchTime", 50, "maximum time to wait for collecting a batch of requests in micro seconds")
	defaultReplica := flag.Int64("defaultReplica", 0, "default replica to send requests to")
	replicaTimeout := flag.Int64("replicaTimeout", 50, "Replica timeout in milli seconds")
	requestSize := flag.Int64("requestSize", 8, "request size in bytes")
	testDuration := flag.Int64("testDuration", 60, "test duration in seconds")
	warmupDuration := flag.Int64("warmupDuration", 100, "warm up duration in seconds")
	arrivalRate := flag.Int64("arrivalRate", 1000, "poisson arrival rate in requests per second")
	benchmark := flag.Int64("benchmark", 0, "Benchmark: 0 for echo service, 1 for KV store and 2 for Redis ")
	numKeys := flag.Int64("numKeys", 1000, "Number of keys in the key value store")
	requestType := flag.String("requestType", "status", "request type: [status , request]")
	operationType := flag.Int64("operationType", 1, "Type of operation for a status request: 1 (bootstrap server), 2: print log)")

	flag.Parse()

	cfg, err := configuration.NewInstanceConfig(*configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}

	cl := cmd.New(*name, cfg, *logFilePath, *batchSize, *batchTime, *defaultReplica, *replicaTimeout, *requestSize, *testDuration, *warmupDuration, *arrivalRate, *benchmark, *numKeys, *requestType, *operationType)

	cl.ConnectToReplicas()
	cl.StartConnectionListners()

}
