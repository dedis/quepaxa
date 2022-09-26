package cmd

import (
	"fmt"
	"github.com/montanaflynn/stats"
	"log"
	"os"
	"raxos/proto"
	"strconv"
	"sync"
)

/*
	iteratively calculate the number of elements in the 2d array
*/
func (cl *Client) getNumberOfSentRequests(requests [][]sentRequestBatch) int {
	count := 0
	for i := 0; i < numRequestGenerationThreads; i++ {
		sentRequestArrayOfI := requests[i] // requests[i] is an array of batch of requests
		for j := 0; j < len(sentRequestArrayOfI); j++ {
			count += len(sentRequestArrayOfI[j].batch.Messages)
		}
	}
	return count
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

	Convert a sync.map to regular map
*/

func (cl *Client) convertToRegularMap(syncMap sync.Map) map[string]receivedResponseBatch {
	var m map[string]receivedResponseBatch
	m = make(map[string]receivedResponseBatch)
	syncMap.Range(func(key, value interface{}) bool {
		m[fmt.Sprint(key)] = value.(receivedResponseBatch)
		return true
	})
	return m
}

/*
	Count the number of individual responses in responses map
*/

func (cl *Client) getNumberOfReceivedResponses(responses sync.Map) int {
	regularMap := cl.convertToRegularMap(responses)

	count := 0
	for _, element := range regularMap {
		count += len(element.batch.Messages)
	}
	return count

}

/*
	Maps the request with the response batch
    Compute the time taken for each request
	Computer the error rate
	Compute the throughput as successfully committed requests per second
	Compute the latency
	Print the basic stats to the stdout
*/

func (cl *Client) computeStats() {

	f, err := os.Create(cl.logFilePath + strconv.Itoa(int(cl.clientName)) + ".txt")
	if err != nil {
		cl.debug("Error creating the output log file", 1)
		log.Fatal(err)
	}
	defer f.Close()

	numTotalSentRequests := cl.getNumberOfSentRequests(cl.sentRequests)
	numTotalReceivedResponses := cl.getNumberOfReceivedResponses(cl.receivedResponses)
	var throughputList []int64 // contains the time duration spent for successful requests
	for i := 0; i < numRequestGenerationThreads; i++ {
		fmt.Printf("Calculating stats for thread %d \n", i)
		for j := 0; j < len(cl.sentRequests[i]); j++ {
			batch := cl.sentRequests[i][j]
			batchId := batch.batch.Id
			matchingResponseBatch, ok := cl.receivedResponses.Load(batchId)
			if ok {
				responseBatch := matchingResponseBatch.(receivedResponseBatch)
				startTime := batch.time
				endTime := responseBatch.time
				batchLatency := endTime.Sub(startTime).Microseconds()
				throughputList = cl.addValueNToArrayMTimes(throughputList, batchLatency, len(batch.batch.Messages))
				cl.printRequests(batch.batch, startTime.Sub(cl.startTime).Microseconds(), endTime.Sub(cl.startTime).Microseconds(), f)
			}

		}
	}

	medianLatency, _ := stats.Median(cl.getFloat64List(throughputList))
	percentile99, _ := stats.Percentile(cl.getFloat64List(throughputList), 99.0) // tail latency
	duration := cl.testDuration
	errorRate := (numTotalSentRequests - len(throughputList)) * 100.0 / numTotalSentRequests

	fmt.Printf("\n Total Sent Requests:= %v \n", numTotalSentRequests)
	fmt.Printf("Total Received Responses:= %v \n", numTotalReceivedResponses)
	fmt.Printf("Total time := %v seconds\n", duration)
	fmt.Printf("Throughput  := %v requests per second\n", len(throughputList)/int(duration))
	fmt.Printf("Median Latency  := %v micro seconds per request\n", medianLatency)
	fmt.Printf("99 pecentile latency  := %v micro seconds per request\n", percentile99)
	fmt.Printf("Error Rate := %v \n", float64(errorRate))
}

/*
	Converts int[] to float64[]
*/

func (cl *Client) getFloat64List(list []int64) []float64 {
	var array []float64
	for i := 0; i < len(list); i++ {
		array = append(array, float64(list[i]))
	}
	return array
}

/*
	Print a client request batch with arrival time and end time
*/

func (cl *Client) printRequests(messages proto.ClientBatch, startTime int64, endTime int64, f *os.File) {
	for i := 0; i < len(messages.Messages); i++ {
		_, _ = f.WriteString(messages.Messages[i].Message + "," + strconv.Itoa(int(startTime)) + "," + strconv.Itoa(int(endTime)) + "\n")
	}
}
