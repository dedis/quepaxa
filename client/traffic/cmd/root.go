package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"github.com/dedis/backsos/config"
	"github.com/spf13/cobra"
)

var (
	err error

	flagConfigFile    string
	numConcurrency    int
	requestPrefix     string
	requestSize       int
	defaultReplica    int
	operationType     int
	arrivalRate       int
	testDuration      int
	warmupDuration    int
	throughputLogFile string
	cfgQuorum         *config.QuorumConfig
	cfgInstances      map[string]string
	workloadType      string
	benchmark         int
	numKeys           int
	requestValue      string
	KVOPs             chan string
	RedisOPs          chan string
	arrival           chan bool
	arrivalTimes      chan int64
	r                 *rand.Rand
	z                 *rand.Zipf
)

func init() {
	cfgInstances = make(map[string]string)
	rootCmd.PersistentFlags().StringVar(&flagConfigFile, "config", "/home/dedis/Baxos/doc/configuration/local/backsosctl/quorum.yml", "Baxos quorum configuration file")
}

var rootCmd = &cobra.Command{
	Use: "backsosctl",
	//	Short: "backsosctl is a control tool for backsos instances",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// load the qorum configuration from file
		cfgQuorum, err = config.NewQuorumConfig(flagConfigFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "load config: %v\n", err)
			os.Exit(1)
		}

		// create a hashmap for easier access to instance addresses by name
		// also define the default instance name (the first one in the list)
		for _, in := range cfgQuorum.Instances {
			cfgInstances[in.Name] = in.Address

		}
	},
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
