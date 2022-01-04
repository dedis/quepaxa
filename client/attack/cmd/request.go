package cmd

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"github.com/dedis/backsos/proto/application"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func getRealSizeOf(v interface{}) (int, error) {
	b := new(bytes.Buffer)
	if err := gob.NewEncoder(b).Encode(v); err != nil {
		return 0, err
	}
	return b.Len(), nil
}

func getStringOfSizeN(length int) string {

	str := "a"
	size, _ := getRealSizeOf(str)
	for size < length {
		str = strings.Repeat(str, 2)
		size, _ = getRealSizeOf(str)
	}

	return str
}

func init() {
	rootCmd.AddCommand(requestCmd)
	requestCmd.PersistentFlags().StringVar(&requestPrefix, "requestPrefix", "requestPrefix", "request prefix")
	requestCmd.PersistentFlags().IntVar(&requestSize, "requestSize", 8, "Size of request")
	requestCmd.PersistentFlags().IntVar(&testDuration, "testDuration", 30, "test duration in seconds")
	requestCmd.PersistentFlags().StringVar(&attackType, "attackType", "delay", "attack type [delay, packet, none]")
	requestCmd.PersistentFlags().StringVar(&dport, "dport", "9000", "[9000, 10000]")
}

var requestCmd = &cobra.Command{
	Use:   "send <request>",
	Short: "Send the request to state machine",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Printf("Starting attacker\n")

		instances := make([]string, 0, len(cfgInstances))

		for k := range cfgInstances {
			instances = append(instances, k)
		}

		numReplicas := len(cfgInstances)

		var clients = make([]application.ApplicationClient, numReplicas)
		var connections = make([]*grpc.ClientConn, numReplicas)

		for w := 0; w < numReplicas; w++ {
			address := cfgInstances[instances[w]]
			//fmt.Printf(" connecting to %v (%v)\n", instances[w], address)
			conn, err := grpc.Dial(address, grpc.WithInsecure())
			if err != nil {
				fmt.Fprintf(os.Stderr, "dial: %v\n", err)
				os.Exit(1)
			}
			connections[w] = conn
			clients[w] = application.NewApplicationClient(conn)
		}

		requestValue := getStringOfSizeN(requestSize)

		fmt.Printf("Finished generating requests\n")
		fmt.Printf("Attack Type: %v\n", attackType)

		start := time.Now()
		for time.Now().Sub(start).Nanoseconds() < int64(testDuration*1000*1000*1000) {
			currentLeader := findCurrentLeader(clients, requestValue, numReplicas)

			if attackType == "delay" {
				attackDelay(cfgInstances[instances[currentLeader]])
			} else if attackType == "packet" {
				attackPacket(cfgInstances[instances[currentLeader]])
			} else {
				time.Sleep(10 * time.Second)
			}

		}

		fmt.Printf("Finished attack: ran for %v\n", time.Now().Sub(start))

		for w := 0; w < numReplicas; w++ {
			_ = connections[w].Close()

		}
	},
}

func findCurrentLeader(clients []application.ApplicationClient, requestValue string, numReplicas int) int {
	currentLeader := rand.Intn(100) % numReplicas
	// connect to instance
retry:
	requestId := strconv.Itoa(1) + requestPrefix

	ctx, cancel := context.WithTimeout(context.Background(), cfgQuorum.Timeout)

	resp, err := clients[currentLeader].Request(ctx, &application.ClientRequest{
		RequestIdentifier: requestId,
		RequestValue:      requestValue,
	})

	if err == nil && resp != nil && resp.ResponseIdentifier == requestId {
		// success
		cancel()
		return currentLeader

	} else {
		currentLeader = (currentLeader + 1) % numReplicas
		cancel()
		goto retry

	}
}

func attackDelay(address string) {

	ip := getIP(address)
	args := []string{"/home/ubuntu/Baxos/delay.sh", ip, dport}
	cmd := exec.Command("/bin/bash", args...)
	err := cmd.Run()
	if err != nil {
		fmt.Printf(" Error while executing delay script %v\n", err)
		log.Fatal(err)
	}

}

func attackPacket(address string) {

	ip := getIP(address)
	args := []string{"/home/ubuntu/Baxos/packet.sh", ip, dport}
	cmd := exec.Command("/bin/bash", args...)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("%v\n", err)
		log.Fatal(err)
	}

}

func getIP(address string) string {
	return strings.Split(address, ":")[0]
}
