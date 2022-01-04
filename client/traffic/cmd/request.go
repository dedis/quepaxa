package cmd

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	benchmark2 "github.com/dedis/backsos/benchmark"
	"github.com/dedis/backsos/proto/application"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"github.com/montanaflynn/stats"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
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
	requestCmd.PersistentFlags().IntVar(&numConcurrency, "numConcurrency", 1000, "number of concurrent threads")
	requestCmd.PersistentFlags().StringVar(&requestPrefix, "requestPrefix", "r", "request prefix")
	requestCmd.PersistentFlags().IntVar(&requestSize, "requestSize", 8, "Size of request in bytes")
	requestCmd.PersistentFlags().IntVar(&testDuration, "testDuration", 60, "test duration in seconds")
	requestCmd.PersistentFlags().IntVar(&warmupDuration, "warmupDuration", 10, "warmup duration in seconds")
	requestCmd.PersistentFlags().IntVar(&defaultReplica, "defaultReplica", 0, "default replica to send the requests to")
	requestCmd.PersistentFlags().IntVar(&arrivalRate, "arrivalRate", 10000, "arrival rate of requests per second")
	requestCmd.PersistentFlags().StringVar(&throughputLogFile, "throughputLogFile", "/home/dedis/Baxos/test.txt", "throughput log file to write the results")
	requestCmd.PersistentFlags().StringVar(&workloadType, "workloadType", "open", "workload type [closed, open]")

	requestCmd.PersistentFlags().IntVar(&benchmark, "benchmark", 0, "0 for no-op, 1 for kv store and 2 for redis")
	requestCmd.PersistentFlags().IntVar(&numKeys, "numKeys", 1000, "number of keys for benchmark")

}

var requestCmd = &cobra.Command{
	Use:   "send <request>",
	Short: "Send the request to state machine",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Printf("Starting client\n")

		wg := sync.WaitGroup{}

		instances := make([]string, 0, len(cfgInstances))

		for k := range cfgInstances {
			instances = append(instances, k)
		}

		numReplicas := len(cfgInstances)

		var clients = make([]application.ApplicationClient, numReplicas)
		var connections = make([]*grpc.ClientConn, numReplicas)

		for w := 0; w < numReplicas; w++ {
			address := cfgInstances[instances[w]]
			conn, err := grpc.Dial(address, grpc.WithInsecure())
			if err != nil {
				fmt.Fprintf(os.Stderr, "dial: %v\n", err)
				os.Exit(1)
			}
			connections[w] = conn
			clients[w] = application.NewApplicationClient(conn)
		}

		type request struct {
			start   int64
			end     int64
			success bool
		}

		requestValue = getStringOfSizeN(requestSize)

		requestsThroughput := make(chan request, 8000000) // change this if the test is run longer than 1 minute

		var testStartTime time.Time
		var testEndTime time.Time

		r = rand.New(rand.NewSource(time.Now().UnixNano()))
		z = rand.NewZipf(r, 3, 10, uint64(numKeys))

		if benchmark == 1 {
			KVOPs = make(chan string, 10000)
			go getKVStoreOp()
			go getKVStoreOp()
		}

		if benchmark == 2 {
			RedisOPs = make(chan string, 10000)
			go getRedisOp()
			go getRedisOp()
		}

		if workloadType == "closed" {

			fmt.Printf("Starting closed loop experiments\n")

			testStartTime = time.Now()

			for i := 0; i < numConcurrency; i++ { // i is level of concurrency

				wg.Add(1)

				go func(k int) {
					defer wg.Done()
					currentLeader := defaultReplica

					start := time.Now()

					j := 0

					for time.Now().Sub(start).Nanoseconds() < int64(testDuration*1000*1000*1000) {
						j++
						// connect to instance
						requestStartTime := time.Now()
						requestId := strconv.Itoa(k) + requestPrefix + strconv.Itoa(j)

						benchmarkRequestValue := ""

						if benchmark == 0 {
							benchmarkRequestValue = requestValue
						} else if benchmark == 1 {
							op := <-KVOPs
							benchmarkRequestValue = op
						} else if benchmark == 2 {
							op := <-RedisOPs
							benchmarkRequestValue = op
						}

						//sending request to currentLeader

						ctx, cancel := context.WithTimeout(context.Background(), cfgQuorum.Timeout)

						resp, err := clients[currentLeader].Request(ctx, &application.ClientRequest{
							RequestIdentifier: requestId,
							RequestValue:      benchmarkRequestValue,
						})

						if err == nil && resp != nil && resp.ResponseIdentifier == requestId {
							// success
							requestEndTime := time.Now()
							requestsThroughput <- request{
								start:   requestStartTime.Sub(start).Microseconds(),
								end:     requestEndTime.Sub(start).Microseconds(),
								success: true,
							}

						} else {
							requestEndTime := time.Now()
							requestsThroughput <- request{
								start:   requestStartTime.Sub(start).Microseconds(),
								end:     requestEndTime.Sub(start).Microseconds(),
								success: false,
							}

							currentLeader = (currentLeader + 1) % numReplicas
						}
						cancel()
					}
				}(i)
			}

			wg.Wait()
			testEndTime = time.Now()
		}

		if workloadType == "open" {

			fmt.Printf("Starting Open Loop Experiments\n")

			arrival = make(chan bool, 20000)
			arrivalTimes = make(chan int64, 20000)

			for i := 0; i < numConcurrency; i++ { // i is level of concurrency

				go func(k int) {
					currentLeader := defaultReplica

					start := time.Now()

					j := 0
					for true { // j is number of requests
						_ = <-arrival
						j++
						// connect to instance
						requestStartTime := time.Now()
						requestId := strconv.Itoa(k) + requestPrefix + strconv.Itoa(j)

						benchmarkRequestValue := ""

						if benchmark == 0 {
							benchmarkRequestValue = requestValue
						} else if benchmark == 1 {
							op := <-KVOPs
							benchmarkRequestValue = op
						} else if benchmark == 2 {
							op := <-RedisOPs
							benchmarkRequestValue = op
						}

						//fmt.Printf("sending request to %v \n", instances[currentLeader])

						ctx, cancel := context.WithTimeout(context.Background(), cfgQuorum.Timeout)
						//fmt.Println("sending request")

						resp, err := clients[currentLeader].Request(ctx, &application.ClientRequest{
							RequestIdentifier: requestId,
							RequestValue:      benchmarkRequestValue,
						})

						if err == nil && resp != nil && resp.ResponseIdentifier == requestId {
							// success
							requestEndTime := time.Now()
							requestsThroughput <- request{
								start:   requestStartTime.Sub(start).Microseconds(),
								end:     requestEndTime.Sub(start).Microseconds(),
								success: true,
							}
						} else {
							requestEndTime := time.Now()
							requestsThroughput <- request{
								start:   requestStartTime.Sub(start).Microseconds(),
								end:     requestEndTime.Sub(start).Microseconds(),
								success: false,
							}

							currentLeader = (currentLeader + 1) % numReplicas
						}
						cancel()

					}

				}(i)
			}

			go generateArrivalTimes()
			// start

			testStartTime = time.Now()

			start := time.Now()

			for time.Now().Sub(start).Nanoseconds() < int64(testDuration*1000*1000*1000) { // j is number if requests

				nextArrivalTime := <-arrivalTimes

				for time.Now().Sub(start).Nanoseconds() < nextArrivalTime {
					// busy waiting
				}
				arrival <- true

			}

			testEndTime = time.Now()
		}

		fmt.Printf("Finished experiment\n")
		time.Sleep(50 * time.Second) // to handle in-flight operations

		close(requestsThroughput)

		ft, errt := os.Create(throughputLogFile)
		if errt != nil {
			panic(errt)
		}

		_, _ = ft.WriteString("Request_Start_Time -- Request_End_Time -- Success \n")

		numLatencyRequests := int64(0)
		var medianLatencyList [] float64
		numSuccess := int64(0)

		for r := range requestsThroughput {
			if r.start > int64(warmupDuration*1000*1000) {
				numLatencyRequests++
				medianLatencyList = append(medianLatencyList, float64(r.end - r.start))
				_, _ = ft.WriteString(strconv.FormatInt(r.start-int64(warmupDuration)*1000*1000, 10) + "--" + strconv.FormatInt(r.end-int64(warmupDuration)*1000*1000, 10) + "--" + strconv.FormatBool(r.success) + "\n")
				if r.success {
					numSuccess++
				}
			}
		}

		_ = ft.Close()

		medianLatency, _ := stats.Median(medianLatencyList)
		percentile99, _ := stats.Percentile(medianLatencyList, 99.0)
		duration := testEndTime.Sub(testStartTime).Seconds() - float64(warmupDuration) // should be changed to reflex the end time of last request - start time of first request
		errorRate := (numLatencyRequests - numSuccess) * 100.0 / numLatencyRequests

		fmt.Printf("Total Requests (includes timeout requests) := %v \n", numLatencyRequests)
		fmt.Printf("Total time := %v seconds\n", duration)
		fmt.Printf("Throughput (successfully committed requests) := %v requests per second\n", float64(numSuccess)/duration)
		fmt.Printf("Median Latency (includes timeout requests) := %v micro seconds per request\n", medianLatency)
		fmt.Printf("99 pecentile latency (includes timeout requests) := %v micro seconds per request\n", percentile99)
		fmt.Printf("Error Rate := %v \n", float64(errorRate))

		for w := 0; w < numReplicas; w++ {
			_ = connections[w].Close()
		}
	},
}

func getKVStoreOp() {
	for true {
		opNum := rand.Intn(2)
		var op string
		if opNum == 0 {
			op = "READ"
		} else {
			op = "UPDATE"
		}

		keyNum := z.Uint64()
		key := benchmark2.Get23LengthRecord(int(keyNum))

		value := ""

		if op == "UPDATE" {
			value = requestValue
		}
		KVOPs <- op + ":" + key + ":" + value

	}
}

func getRedisOp() {
	for true {
		opNum := rand.Intn(2)
		var op string

		if opNum == 0 {
			op = "READ"
		} else {
			op = "UPDATE"
		}

		keyNum := z.Uint64()
		key := benchmark2.Get23LengthRecord(int(keyNum))

		field := "field0"

		value := ""
		if op == "UPDATE" {
			value = benchmark2.Get100LengthValue()
		}
		RedisOPs <- op + ":" + key + ":" + field + ":" + value

	}

}

func generateArrivalTimes() {
	lambda := float64(arrivalRate) / (1000.0 * 1000.0 * 1000.0) // requests per nano second
	arrivalTime := 0.0

	for true {
		// Get the next probability value from Uniform(0,1)
		p := rand.Float64()

		//Plug it into the inverse of the CDF of Exponential(_lamnbda)
		interArrivalTime := -1 * (math.Log(1.0-p) / lambda)

		// Add the inter-arrival time to the running sum
		arrivalTime = arrivalTime + interArrivalTime

		arrivalTimes <- int64(arrivalTime)
	}
}
