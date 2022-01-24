package cmd

import (
	"math"
	"math/rand"
	"raxos/proto"
	raxos "raxos/replica/src"
	"strconv"
	"time"
)

/*
	Upon receiving a client response, add the request to the received requests array
*/

func (cl *Client) handleClientResponseBatch(batch *proto.ClientResponseBatch) {
	cl.receivedResponses = append(cl.receivedResponses, receivedResponseBatch{
		batch: *batch,
		time:  time.Now(), // record the time when the response was received
	})
	cl.lastSeenTimeReplica = time.Now() // mark the last time a response was received
	cl.debug("Received Response Batch")
}

/*
	start the poisson arrival process (put arrivals to arrivalTimeChan) in a separate thread
	start request generation processes  (get arrivals from arrivalTimeChan and generate batches and send them) in separate threads, and send them to the leader proposer, and write the sent batch to the correct array in sentRequests. Runs only for a pre-defined test duration
	start failure detector that checks the time since the last response was received, and update the default proposer
	the thread sleeps for test duration + delta and then starts processing the responses

*/

func (cl *Client) SendRequests() {
	cl.generateArrivalTimes()
	cl.startRequestGenerators()
	cl.startFailureDetector()
	cl.startScheduler() // this is sync, main thread waits for this to finish

	// end of test

	time.Sleep(time.Duration(cl.testDuration) * time.Second) // additional sleep duration to make sure that all the inflight responses are received
	cl.computeStats()
}

/*
	Each request generator generates requests by generating string requests, forming batches and add them to the correct sent array
*/

func (cl *Client) startRequestGenerators() {
	for i := 0; i < numRequestGenerationThreads; i++ { // i is the thread number
		go func(threadNumber int) {
			localCounter := 0
			lastSent := time.Now() // used to get how long to wait
			for true {             // this runs forever
				numRequests := int64(0)
				sampleRequest := "Request" //todo this is only for testing purposes
				var requests []*proto.ClientRequestBatch_SingleClientRequest
				// this loop collects requests until the minimum batch size is met OR the batch time is timeout
				for !(numRequests >= cl.batchSize || (time.Now().Sub(lastSent).Microseconds() > cl.batchTime && numRequests > 0)) {
					_ = <-cl.arrivalChan                                                                               // keep collecting new requests arrivals
					requests = append(requests, &proto.ClientRequestBatch_SingleClientRequest{Message: sampleRequest}) //todo for actual benchmarks the same request should be replaced with redis op or kvstore op
					numRequests++
				}

				batch := proto.ClientRequestBatch{
					Sender:   cl.clientName,
					Receiver: cl.defaultReplica,
					Requests: requests,
					Id:       strconv.Itoa(int(cl.clientName)) + "." + strconv.Itoa(threadNumber) + "." + strconv.Itoa(localCounter), // this is a unique string
				}

				cl.debug("Sent " + strconv.Itoa(int(cl.clientName)) + "." + strconv.Itoa(threadNumber) + "." + strconv.Itoa(localCounter) + " batch size " + strconv.Itoa(len(requests)))

				localCounter++

				rpcPair := raxos.RPCPair{
					Code: cl.clientRequestBatchRpc,
					Obj:  &batch,
				}

				cl.sendMessage(cl.defaultReplica, rpcPair)
				lastSent = time.Now()
				cl.sentRequests[threadNumber] = append(cl.sentRequests[threadNumber], sentRequestBatch{
					batch: batch,
					time:  time.Now(),
				})
			}

		}(i)
	}

}

/*
	Until the test duration is arrived, fetch new arrivals and inform the request generators

*/

func (cl *Client) startScheduler() {
	start := time.Now()

	for time.Now().Sub(start).Nanoseconds() < cl.testDuration*1000*1000*1000 {
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
	go func() {
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
	}()
}

/*
	Monitors the time the last response was received. If the default replica fails to send a response before a timeout, change the default replica
	Since the failure detector is anyway eventually correct, we don't use a Mutex to protect the lastSeenTime and the defaultReplica variables
*/

func (cl *Client) startFailureDetector() {
	go func() {
		cl.debug("Starting failure detector")
		for true {

			time.Sleep(time.Duration(cl.replicaTimeout) * time.Second)
			if time.Now().Sub(cl.lastSeenTimeReplica).Seconds() > float64(cl.replicaTimeout) {

				// change the default replica
				cl.debug("Changing the default replica")
				cl.defaultReplica = (cl.defaultReplica + 1) % cl.numReplicas
				cl.lastSeenTimeReplica = time.Now()

			}

		}
	}()
}
