package main

import (
	"flag"
	"fmt"
	"os"
	"raxos/client/cmd"
	"raxos/configuration"
	"time"
)

/*
	This file defines the main function of the client.
	The main function performs the following steps:
	Collects the command line input arguments. Descriptions for each command line argument can be found in the "usage" field.
	Parses the configuration file, which is automatically generated by the scripts.
	Provides two client methods:
	a. Request client: Sends a stream of requests in a Poisson open loop, with a specified window to allow back pressure.
	b. Status client: Sends status operations, which are further explained in client/cmd/status.go.
*/

func main() {
	name := flag.Int64("name", 21, "name of the client [ 21, 21+numClients) ]")
	configFile := flag.String("config", "configuration/local/configuration.yml", "raxos configuration file which contains the ip:port of each client and replica, this file is auto generated")
	logFilePath := flag.String("logFilePath", "logs/", "log file path")
	batchSize := flag.Int64("batchSize", 50, "client batch size")
	batchTime := flag.Int64("batchTime", 50, "client batch time in micro seconds")
	testDuration := flag.Int64("testDuration", 60, "test duration in seconds")
	arrivalRate := flag.Int64("arrivalRate", 10000, "Poisson arrival rate in requests per second")
	requestType := flag.String("requestType", "status", "request type: [status , request]")
	operationType := flag.Int64("operationType", 1, "Type of operation for a status request: 1 bootstrap server, 2: print log, 3: slow down proposal speed for simulation, and 4. average number of steps per slot ")
	keyLen := flag.Int64("keyLen", 8, "length of key in client requests")
	valLen := flag.Int64("valLen", 8, "length of value in client requests")
	debugOn := flag.Bool("debugOn", false, "turn on/off debug, turn off when benchmarking")
	debugLevel := flag.Int64("debugLevel", 0, "debug level, debug messages with equal or higher debugLevel will be printed")
	slowdown := flag.String("slowdown", "", "node1:wait1,node2:wait2,node3:wait3  artificial delay for each node id, in micro seconds ")
	window := flag.Int64("window", 1000, "number of outstanding client batches sent by the client, before receiving the response")

	flag.Parse()
	// pass the config file ip addresses into a config object cfg
	cfg, err := configuration.NewInstanceConfig(*configFile, *name)

	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}
	// create a new client instance
	cl := cmd.New(*name, cfg, *logFilePath, *batchSize, *batchTime, *testDuration, *arrivalRate, *requestType, *operationType, int(*debugLevel), *debugOn, int(*keyLen), int(*valLen), *slowdown, *window)

	// start the threads which receive out going messages
	cl.StartOutgoingLinks()

	// start listening to all incoming replica connections
	go cl.WaitForConnections()

	// main single threaded event loop of client
	cl.Run()

	// connect to all replicas
	cl.ConnectToReplicas()

	time.Sleep(5 * time.Second) // this wait time is for the connections to establish

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
