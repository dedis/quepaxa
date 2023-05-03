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
	batchTime := flag.Int64("batchTime", 5000, "replica batch time in micro seconds")
	leaderTimeout := flag.Int64("leaderTimeout", 5000000, "leader timeout in micro seconds")
	pipelineLength := flag.Int64("pipelineLength", 50, "pipeline length maximum number of outstanding proposals")
	benchmark := flag.Int64("benchmark", 0, "Benchmark: 0 for KV store and 1 for Redis ")
	debugOn := flag.Bool("debugOn", false, "true / false")
	debugLevel := flag.Int("debugLevel", 0, "debug level")
	leaderMode := flag.Int("leaderMode", 0, "mode of leader change: 0 for fixed leader order, 1 for fixed order, static partition,  2 for M.A.B based on commit times, 3 for asynchronous")
	serverMode := flag.Int("serverMode", 0, "0 for non-lan-optimized, 1 for lan optimized")
	epochSize := flag.Int("epochSize", 100, "epoch size for MAB")
	keyLen := flag.Int64("keyLen", 8, "length of key")
	valLen := flag.Int64("valLen", 8, "length of value")
	requestPropagationTime := flag.Int64("requestPropagationTime", 0, "additional wait time in 'milli seconds' for client batches, such that there is enough time for client driven request propagation, for LAN this is 0, for WAN it is usually set to 5ms")
	checkProposerDuplicates := flag.Bool("checkProposerDuplicates", false, "false if no checking is required, turn on for asynchronous experiments")
	flag.Parse()

	cfg, err := configuration.NewInstanceConfig(*configFile, *name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}

	serverInstance := raxos.New(cfg, *name, *logFilePath, *batchSize, *leaderTimeout, *pipelineLength, *debugOn, *debugLevel, *leaderMode, *serverMode, *batchTime, *epochSize, int(*benchmark), int(*keyLen), int(*valLen), *checkProposerDuplicates, *requestPropagationTime) // create a new server instance

	serverInstance.NetworkInit()
	serverInstance.Run()

	/*to avoid exiting the main thread*/
	for true {
		time.Sleep(10 * time.Second)
	}
}
