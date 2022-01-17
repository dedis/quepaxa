package cmd

import (
	"math"
	"math/rand"
	"raxos/proto"
	"time"
)

/*
Upon receiving a client response, add the request to the received requests array
*/

func (cl *Client) handleClientResponseBatch(batch *proto.ClientResponseBatch) {
	cl.receivedResponses = append(cl.receivedResponses, receivedResponseBatch{
		batch: *batch,
		time:  time.Now(),
	})
}

/*
	start the poisson arrival process (put arrivals to arrivalTimeChan) in a separate thread
	start request generation processes  (get arrivals from arrivalTimeChan and generate batches and send them) in separate threads, and send them to the leader proposer, and write the sent batch to the correct array in sentRequests. Runs only for a pre-defined test duration
	start failure detector that checks the time since the last response was received, and update the default proposer
	the main thread sleeps for test duration + delta and then starts processing the responses

*/

func (cl *Client) SendRequests() {
	go cl.generateArrivalTimes()

	cl.startScheduler() // this is sync, main thread waits for this to finish

}

/*
	Each request generator generates requests by generating string requests, forming batches and add them to the correct sent array
*/

func (cl *Client) startRequestGenerators() {
	for i := 0; i < numRequestGenerationThreads; i++ {

	}
}

/*
	Until the test duration is arrived, fetch new arrivals and inform the request generators

*/

func (cl *Client) startScheduler() {
	start := time.Now()

	for time.Now().Sub(start).Nanoseconds() < int64(cl.testDuration*1000*1000*1000) {
		nextArrivalTime := <-cl.arrivalTimeChan

		for time.Now().Sub(start).Nanoseconds() < nextArrivalTime {
			// busy waiting until the time to dispatch this request arrives
		}
		cl.arrivalChan <- true

	}
}

/*
 Generates Poisson arrival times
*/

func (cl *Client) generateArrivalTimes() {
	lambda := float64(cl.arrivalRate) / (1000.0 * 1000.0 * 1000.0) // requests per nano second
	arrivalTime := 0.0

	for true {
		// Get the next probability value from Uniform(0,1)
		p := rand.Float64()

		//Plug it into the inverse of the CDF of Exponential(_lamnbda)
		interArrivalTime := -1 * (math.Log(1.0-p) / lambda)

		// Add the inter-arrival time to the running sum
		arrivalTime = arrivalTime + interArrivalTime

		cl.arrivalTimeChan <- int64(arrivalTime)
	}
}
