package main

import (
	"flag"
	"fmt"
	"os"
	"raxos/client/cmd"
	"raxos/configuration"
	"time"
)

func main() {
	name := flag.Int64("name", 5, "name of the client as specified in the configuration.yml")
	configFile := flag.String("config", "configuration/local/configuration.yml", "raxos configuration file")
	logFilePath := flag.String("logFilePath", "logs/", "log file path")
	batchSize := flag.Int64("batchSize", 50, "client batch size")
	requestSize := flag.Int64("requestSize", 8, "request size in bytes")
	testDuration := flag.Int64("testDuration", 60, "test duration in seconds")
	arrivalRate := flag.Int64("arrivalRate", 10000, "poisson arrival rate in requests per second")
	benchmark := flag.Int64("benchmark", 0, "Benchmark: 0 for echo service, 1 for KV store and 2 for Redis ")
	requestType := flag.String("requestType", "status", "request type: [status , request]")
	operationType := flag.Int64("operationType", 1, "Type of operation for a status request: 1 (bootstrap server, 2: print log)")
	debugLevel := flag.Int64("debugLevel", 0, "debug level")
	debugOn := flag.Bool("debugOn", false, "turn on/off debug")

	flag.Parse()

	cfg, err := configuration.NewInstanceConfig(*configFile, *name)

	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}

	cl := cmd.New(*name, cfg, *logFilePath, *batchSize, *requestSize, *testDuration, *arrivalRate, *benchmark, *requestType, *operationType, int(*debugLevel), *debugOn)

	cl.StartOutgoingLinks()
	go cl.WaitForConnections()
	cl.Run()
	cl.ConnectToReplicas()

	time.Sleep(5 * time.Second)

	if cl.RequestType == "status" {
		cl.SendStatus(cl.OperationType)
	} else if cl.RequestType == "request" {
		cl.SendRequests()
	}

}
