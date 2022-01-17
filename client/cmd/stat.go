package cmd

import (
	"fmt"
	"github.com/montanaflynn/stats"
)

/*
	iteratively calculates the number of elements in the 2d array
*/
func (cl *Client) getNumberofRequests(requests [][]sentRequestBatch) int {
	count := 0
	for i := 0; i < numRequestGenerationThreads; i++ {
		count += len(requests[i])
	}
	return count
}

/*
	returns the matching response batch corresponding to the request batch id
*/

func (cl *Client) getMatchingResponseBatch(id string) int {
	for i := 0; i < len(cl.receivedResponses); i++ {
		if cl.receivedResponses[i].batch.Id == id {
			return i
		}
	}
	return -1
}

/*
	Add value N to list, M times
*/

func (cl *Client) addValueNToArrayMTimes(list []int64, N int64, M int) []int64 {
	for i := 0; i < M; i++ {
		list = append(list, N)
	}
	return list
}

/*
	Maps the request with the response batch
    Compute the time taken for each request
	Computer the error rate
	Compute the throughput as successfully committed requests per second (doesn't include failued requests)
	Compute the latency by including the timeout requests
	Prints the basic stats to the stdout
*/

func (cl *Client) computeStats() {

	numTotalRequests := cl.getNumberofRequests(cl.sentRequests)
	var latencyList []int64    // contains the time duration spent for each request
	var throughputList []int64 // contains the time duration spent for successful requests
	for i := 0; i < numRequestGenerationThreads; i++ {
		for j := 0; j < len(cl.sentRequests[i]); j++ {
			batch := cl.sentRequests[i][j]
			batchId := batch.batch.Id
			matchingResponseIndex := cl.getMatchingResponseBatch(batchId)
			if matchingResponseIndex == -1 {
				// there is no response for this batch of requests, hence they are considered as failed
				cl.addValueNToArrayMTimes(latencyList, cl.replicaTimeout, len(batch.batch.Requests))
			} else {
				responseBatch := cl.receivedResponses[matchingResponseIndex]
				startTime := batch.time
				endTime := responseBatch.time
				bacthLatency := endTime.Sub(startTime).Microseconds()
				cl.addValueNToArrayMTimes(latencyList, bacthLatency, len(batch.batch.Requests))
				cl.addValueNToArrayMTimes(throughputList, bacthLatency, len(batch.batch.Requests))
			}

		}
	}

	medianLatency, _ := stats.Median(cl.getFloat64List(latencyList))
	percentile99, _ := stats.Percentile(cl.getFloat64List(latencyList), 99.0)
	duration := cl.testDuration
	errorRate := (numTotalRequests - len(throughputList)) * 100.0 / numTotalRequests

	fmt.Printf("Total Requests (includes timeout requests) := %v \n", numTotalRequests)
	fmt.Printf("Total time := %v seconds\n", duration)
	fmt.Printf("Throughput (successfully committed requests) := %v requests per second\n", len(throughputList)/int(duration))
	fmt.Printf("Median Latency (includes timeout requests) := %v micro seconds per request\n", medianLatency)
	fmt.Printf("99 pecentile latency (includes timeout requests) := %v micro seconds per request\n", percentile99)
	fmt.Printf("Error Rate := %v \n", float64(errorRate))
}

func (cl *Client) getFloat64List(list []int64) []float64 {
	var array []float64
	for i := 0; i < len(list); i++ {
		array = append(array, float64(list[i]))
	}
	return array
}
