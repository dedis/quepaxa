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
	name := flag.Int64("name", 21, "name of the client as specified in the configuration.yml")
	configFile := flag.String("config", "configuration/local/configuration.yml", "raxos configuration file")
	logFilePath := flag.String("logFilePath", "logs/", "log file path")
	batchSize := flag.Int64("batchSize", 50, "client batch size")
	batchTime := flag.Int64("batchTime", 50, "client batch time micro seconds")
	testDuration := flag.Int64("testDuration", 60, "test duration in seconds")
	arrivalRate := flag.Int64("arrivalRate", 10000, "poisson arrival rate in requests per second")
	requestType := flag.String("requestType", "status", "request type: [status , request]")
	operationType := flag.Int64("operationType", 1, "Type of operation for a status request: 1 bootstrap server, 2: print log, 3: slow down proposal speed")
	keyLen := flag.Int64("keyLen", 8, "length of key")
	valLen := flag.Int64("valLen", 8, "length of value")
	debugOn := flag.Bool("debugOn", false, "turn on/off debug")
	debugLevel := flag.Int64("debugLevel", 0, "debug level")
	slowdown := flag.String("slowdown", "", "node1:wait1,node2:wait2,node3:wait3")
	window := flag.Int64("window", 1000, "number of outstanding client batches")

	flag.Parse()

	cfg, err := configuration.NewInstanceConfig(*configFile, *name)

	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}

	cl := cmd.New(*name, cfg, *logFilePath, *batchSize, *batchTime, *testDuration, *arrivalRate, *requestType, *operationType, int(*debugLevel), *debugOn, int(*keyLen), int(*valLen), *slowdown, *window)

	cl.StartOutgoingLinks()
	go cl.WaitForConnections()
	cl.Run()
	cl.ConnectToReplicas()

	time.Sleep(5 * time.Second)

	if cl.RequestType == "status" {
		fmt.Printf("\nstarting status client\n")
		cl.SendStatus(cl.OperationType)
		fmt.Printf("\nfinishing status client\n")
	} else if cl.RequestType == "request" {
		fmt.Printf("\nstarting request client\n")
		cl.SendRequests()
		fmt.Printf("\nfinishing request client\n")
	}

}
