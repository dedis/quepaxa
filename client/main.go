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
	testDuration := flag.Int64("testDuration", 60, "test duration in seconds")
	arrivalRate := flag.Int64("arrivalRate", 10000, "poisson arrival rate in requests per second")
	requestType := flag.String("requestType", "status", "request type: [status , request]")
	operationType := flag.Int64("operationType", 1, "Type of operation for a status request: 1 (bootstrap server, 2: print log)")
	keyLen := flag.Int64("keyLen", 8, "length of key")
	valLen := flag.Int64("valLen", 8, "length of value")
	debugOn := flag.Bool("debugOn", false, "turn on/off debug")
	debugLevel := flag.Int64("debugLevel", 0, "debug level")

	flag.Parse()

	cfg, err := configuration.NewInstanceConfig(*configFile, *name)

	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}

	cl := cmd.New(*name, cfg, *logFilePath, *batchSize, *testDuration, *arrivalRate, *requestType, *operationType, int(*debugLevel), *debugOn, int(*keyLen), int(*valLen))

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
