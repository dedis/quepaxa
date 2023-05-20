package cmd

import (
	"fmt"
	"math"
	"math/rand"
	"raxos/common"
	"raxos/proto/client"
	"strconv"
	"time"
	"unsafe"
)

/*
	this file defines the client request and response logics
	request sending sends a stream of batched client requests in a Poisson open loop with back pressure
	response handling just saves the reveived response
*/

/*
	Upon receiving a client response batch, add the batch to the received requests map
*/

func (cl *Client) handleClientResponseBatch(batch *client.ClientBatch) {

	// we are done, then do not handle response
	if cl.finished {
		return
	}
	// if we already received this response batch, then ignore
	_, ok := cl.receivedResponses.Load(batch.Id)
	if ok {
		return
	}

	cl.receivedResponses.Store(batch.Id, receivedResponseBatch{
		batch: *batch,
		time:  time.Now(), // record the time when the response was received
	})
	cl.receiveCountMutex.Lock()
	cl.totalReceivedBatches++
	cl.receiveCountMutex.Unlock()
	//cl.debug("Added response Batch from "+strconv.Itoa(int(batch.Sender))+" to received map", 0)
}

/*
	start the poisson arrival process (put arrivals to arrivalTimeChan) in a separate thread
	start request generation processes  (get arrivals from arrivalTimeChan and generate batches and send them) in separate threads,
and send them to all the proxies, and write batch to the correct array in sentRequests
	start the scheduler that schedules new requests
	the thread sleeps for test duration + delta and then starts processing the responses
*/

func (cl *Client) SendRequests() {
	cl.generateArrivalTimes()
	cl.startRequestGenerators()
	cl.startScheduler() // this is sync, main thread waits for this to finish

	// end of test

	time.Sleep(time.Duration(20) * time.Second) // additional sleep duration to make sure that all the in-flight responses are received
	fmt.Printf("Finish sending requests \n")
	cl.finished = true
	cl.computeStats()
}

/*
	random string generation adapted from the Rabia SOSP 2021 code base https://github.com/haochenpan/rabia/
*/

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // low conflict
	letterIdxBits = 6                                                      // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1                                   // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits                                     // # of letter indices fitting in 63 bits
)

/*
	generate a random string of length n adapted from the Rabia SOSP 2021 code base https://github.com/haochenpan/rabia/
*/

func (cl *Client) RandString(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

/*
	request generator generates requests by generating string requests, forming batches, send batches and add them to the correct sent array
*/

func (cl *Client) startRequestGenerators() {

	for i := 0; i < numRequestGenerationThreads; i++ { // i is the thread number
		go func(threadNumber int) {
			localCounter := 0
			lastSent := time.Now()
			for true {
				if cl.finished {
					return
				}
				numRequests := int64(0)
				var requests []*client.ClientBatch_SingleMessage
				// this loop collects requests until the minimum batch size or batch time is met
				for !(numRequests >= cl.batchSize || (time.Now().Sub(lastSent).Microseconds() > cl.batchTime && numRequests > 0)) {
					_ = <-cl.arrivalChan // keep collecting new requests arrivals
					requests = append(requests, &client.ClientBatch_SingleMessage{
						Message: fmt.Sprintf("%d%v%v", rand.Intn(2),
							cl.RandString(cl.keyLen),
							cl.RandString(cl.valLen)),
					})
					numRequests++
				}
				cl.receiveCountMutex.Lock()
				if (cl.totalSentBatches - cl.totalReceivedBatches) > cl.window {
					cl.receiveCountMutex.Unlock()
					continue
				}
				cl.receiveCountMutex.Unlock()

				for i, _ := range cl.replicaAddrList {

					var requests_i []*client.ClientBatch_SingleMessage

					for j := 0; j < len(requests); j++ {
						requests_i = append(requests_i, requests[j])
					}

					batch := client.ClientBatch{
						Sender:   cl.name,
						Messages: requests_i,
						Id:       strconv.Itoa(int(cl.name)) + "." + strconv.Itoa(threadNumber) + "." + strconv.Itoa(localCounter), // this is a unique string id,
					}

					rpcPair := common.RPCPair{
						Code: cl.clientBatchRpc,
						Obj:  &batch,
					}

					cl.sendMessage(i, rpcPair)
				}

				//cl.debug("Sent "+strconv.Itoa(int(cl.name))+"."+strconv.Itoa(threadNumber)+"."+strconv.Itoa(localCounter), 0)
				cl.totalSentBatches++

				batch := client.ClientBatch{
					Sender:   cl.name,
					Messages: requests,
					Id:       strconv.Itoa(int(cl.name)) + "." + strconv.Itoa(threadNumber) + "." + strconv.Itoa(localCounter), // this is a unique string id,
				}

				cl.sentRequests[threadNumber] = append(cl.sentRequests[threadNumber], sentRequestBatch{
					batch: batch,
					time:  time.Now(),
				})
				lastSent = time.Now()
				localCounter++
			}
		}(i)
	}

}

/*
	until the test duration is arrived, fetch new arrivals and inform the request generators
*/

func (cl *Client) startScheduler() {

	cl.startTime = time.Now()

	for time.Now().Sub(cl.startTime).Nanoseconds() < cl.testDuration*1000*1000*1000 {
		nextArrivalTime := <-cl.arrivalTimeChan

		for time.Now().Sub(cl.startTime).Nanoseconds() < nextArrivalTime {
			// busy waiting until the time to dispatch this request arrives
		}
		cl.arrivalChan <- true
	}
}

/*
	generate poisson arrival times
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
