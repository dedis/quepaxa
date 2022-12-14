package raxos

import (
	"fmt"
	"time"
)

// checks if all indexes in this epoch are decided

func (pr *Proxy) hasAllDecided(epoch int) bool {
	startIndex := epoch * pr.epochSize
	endIndex := (epoch+1)*pr.epochSize - 1

	for i := startIndex; i <= endIndex; i++ {
		if !(len(pr.replicatedLog) >= i+1 && pr.replicatedLog[i].decided == true) {
			return false
		}
	}
	return true
}

// update the start time and the end time of the epoch

func (pr *Proxy) updateEpochTime(index int) {
	epoch := index / pr.epochSize
	for len(pr.epochTimes) < epoch+1 {
		pr.epochTimes = append(pr.epochTimes, EpochTime{
			startTime: time.Time{},
			endTime:   time.Time{},
			started:   false,
			ended:     false,
		})
	}

	if pr.epochTimes[epoch].started == false {
		pr.debug("starting epoch "+fmt.Sprintf("%v at time %v", epoch, time.Now()), 9)
		pr.epochTimes[epoch].started = true
		pr.epochTimes[epoch].startTime = time.Now()
	}

	if pr.epochTimes[epoch].ended == false && pr.hasAllDecided(epoch) {
		pr.debug("finishing epoch "+fmt.Sprintf("%v at time %v", epoch, time.Now()), 9)
		pr.epochTimes[epoch].ended = true
		pr.epochTimes[epoch].endTime = time.Now()
	}
}

// return the immutable leader sequence for instance

func (pr *Proxy) getLeaderSequence(instance int64) []int64 {
	if pr.leaderMode == 0 {
		// fixed order
		// assumes that node names start with 1
		rA := make([]int64, 0)

		for i := 0; i < pr.numReplicas; i++ {
			rA = append(rA, int64(i+1))
		}

		return rA
	}
	if pr.leaderMode == 1 {
		// fixed order, static partition
		// assumes that node names start with 1
		epoch := instance / int64(pr.epochSize)
		sequence := epoch % int64(pr.numReplicas) // sequence is 0-numreplicas

		rA := make([]int64, 0)
		for i := sequence; i < int64(pr.numReplicas); i++ {
			rA = append(rA, i+1)
		}
		for i := int64(0); i < sequence; i++ {
			rA = append(rA, i+1)
		}
		pr.debug("proxy leader sequence for instance "+fmt.Sprintf("%v is %v", instance, rA), 0)
		return rA

	}
	// todo

	panic("should not happen")
}

// return the pre-agreed, non changing waiting time for the instance by the proposer todo

func (pr *Proxy) getLeaderWait(sequence []int64) int64 {

	for j := 0; j < len(sequence); j++ {
		if sequence[j] == pr.name {
			return pr.leaderTimeout * int64(j)
		}
	}
	panic("should not happen")
}
