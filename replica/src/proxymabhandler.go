package raxos

import (
	"fmt"
	"raxos/proto/client"
	"sort"
	"strconv"
	"strings"
	"time"
)

// checks if all indexes in this epoch are decided

func (pr *Proxy) hasAllDecided(epoch int) bool {
	startIndex := epoch * pr.epochSize
	endIndex := (epoch+1)*pr.epochSize - 2 // last index in the epoch is a bit biased

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
		pr.debug("starting epoch "+fmt.Sprintf("%v at time %v", epoch, time.Now()), 11)
		pr.epochTimes[epoch].started = true
		pr.epochTimes[epoch].startTime = time.Now()
	}

	if pr.epochTimes[epoch].ended == false && pr.hasAllDecided(epoch) {
		pr.debug("finishing epoch "+fmt.Sprintf("%v at time %v", epoch, time.Now()), 11)
		pr.epochTimes[epoch].ended = true
		pr.epochTimes[epoch].endTime = time.Now()
		pr.debug("epoch "+fmt.Sprintf("%v took %v ms", epoch, pr.epochTimes[epoch].endTime.Sub(pr.epochTimes[epoch].startTime).Milliseconds()), 10)
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
	if pr.leaderMode == 2 {
		// M.A.B based on commit times for each epoch
		// assumes that node names start with 1
		epoch := instance / int64(pr.epochSize)

		if epoch <= 2*int64(pr.numReplicas) {
			// first 2 x numReplicas +1 epochs we do round-robin
			sequence := epoch % int64(pr.numReplicas) // sequence is 0-numreplicas

			rA := make([]int64, 0)
			for i := sequence; i < int64(pr.numReplicas); i++ {
				rA = append(rA, i+1)
			}
			for i := int64(0); i < sequence; i++ {
				rA = append(rA, i+1)
			}
			pr.debug("proxy leader sequence for instance "+fmt.Sprintf("%v is %v", instance, rA), 10)
			return rA
		} else {
			rA := pr.getLeaderSequenceFromLastEpoch(epoch)
			pr.debug("proxy leader sequence for instance "+fmt.Sprintf("%v is %v", instance, rA), 10)
			return rA
		}

	}

	panic("should not happen")
}

// return the pre-agreed, non changing waiting time for the instance by the proposer

func (pr *Proxy) getLeaderWait(sequence []int64) int64 {

	for j := 0; j < len(sequence); j++ {
		if sequence[j] == pr.name {
			return pr.leaderTimeout * int64(j)
		}
	}
	panic("should not happen")
}

// return the epoch - 1 th epoch's first element

func (pr *Proxy) getLeaderSequenceFromLastEpoch(epoch int64) []int64 {
	index := (epoch - 1) * int64(pr.epochSize)
	if pr.replicatedLog[index].decided != true {
		panic("should this happen")
	}
	decision := pr.replicatedLog[index].decidedBatch[0]
	return pr.convertToIntArray(decision)
}

// the string is in the form of Epoch1,2,3,4,5 convert it to int64[]

func (pr *Proxy) convertToIntArray(decisionId string) []int64 {
	// decision id is int form Epoch1,2,3,4,5
	s := strings.Split(decisionId[5:], ",")
	intArr := make([]int64, pr.numReplicas)
	if len(s) != pr.numReplicas {
		panic("should this happen?")
	}
	for i := 0; i < pr.numReplicas; i++ {
		num, _ := strconv.Atoi(s[i])
		intArr[i] = int64(num)
	}
	return intArr
}

// checks if the corresponding epoch is greater than numReplicas and index is the first element of the epoch
func (pr *Proxy) isBeginningOfEpoch(index int64) bool {
	epoch := index / int64(pr.epochSize)

	if epoch >= int64(2*pr.numReplicas) {
		startIndex := epoch * int64(pr.epochSize)
		if index == startIndex {
			return true
		}
	}
	return false
}

// propose the leader sequence for the 1st instance of each epoch

func (pr *Proxy) proposePreviousEpochSummary(index int64) {
	curEpoch := index / int64(pr.epochSize)
	var strSequence string
	strSequence = pr.calculateSequence(int(curEpoch))

	strProposals := make([]string, 0)
	btchProposals := make([]client.ClientBatch, 0)

	strProposals = []string{"Epoch" + strSequence}
	btchProposals = append(btchProposals, client.ClientBatch{
		Sender:   -1,
		Messages: []*client.ClientBatch_SingleMessage{{Message: "Epoch" + strSequence}},
		Id:       "Epoch" + strSequence,
	})
	pr.debug("proposing new summary for index "+fmt.Sprintf("%v, epoch:%v, sequence:%v ", index, curEpoch, strSequence), 13)

	newProposalRequest := ProposeRequest{
		instance:             index,
		proposalStr:          strProposals,
		proposalBtch:         btchProposals,
		msWait:               int(pr.getLeaderWait(pr.getLeaderSequence(index))),
		lastDecidedIndexes:   pr.lastDecidedIndexes,
		lastDecidedDecisions: pr.lastDecidedDecisions,
		leaderSequence:       pr.getLeaderSequence(index),
	}

	pr.proxyToProposerChan <- newProposalRequest
	pr.debug("proxy sent a proposal request containing leader sequence  "+fmt.Sprintf("%v", newProposalRequest), 10)
	// create the slot index
	for len(pr.replicatedLog) < int(index+1) {
		// create the new entry
		pr.replicatedLog = append(pr.replicatedLog, Slot{
			proposedBatch: nil,
			decidedBatch:  nil,
			decided:       false,
			committed:     false,
		})
	}

	pr.replicatedLog[index] = Slot{
		proposedBatch: strProposals,
		decidedBatch:  pr.replicatedLog[index].decidedBatch,
		decided:       pr.replicatedLog[index].decided,
		committed:     pr.replicatedLog[index].committed,
	}

	// reset the variables
	pr.lastDecidedIndexes = make([]int, 0)
	pr.lastDecidedDecisions = make([][]string, 0)
}

// statistically calculate the leader sequence in the form 1,2,3,4,5

func (pr *Proxy) calculateSequence(epoch int) string {

	if epoch < 2*pr.numReplicas {
		panic("should not happen")
	}

	times := make([][]int64, pr.numReplicas)

	for i := 0; i < pr.numReplicas; i++ {
		times[i] = make([]int64, 0)
	}

	for i := 0; i < epoch-1; i++ {

		if pr.epochTimes[i].ended != true {
			panic("should this happen")
		}
		ld := pr.getLeaderSequence(int64(i * pr.epochSize))[0]
		times[ld-1] = append(times[ld-1], pr.epochTimes[i].endTime.Sub(pr.epochTimes[i].startTime).Milliseconds())
	}

	pr.debug("epoch time summary "+fmt.Sprintf("for the epoch %v is %v ", epoch, times), 10)

	epochTimes1 := make([]int, pr.numReplicas)
	epochTimes2 := make([]int, pr.numReplicas)

	for i := 0; i < pr.numReplicas; i++ {
		sum := int64(0)
		count := int64(0)
		for j := 0; j < len(times[i]); j++ {
			sum += times[i][j]
			count++
		}
		epochTimes1[i] = int(sum / count)
		epochTimes2[i] = int(sum / count)
	}
	pr.debug("epoch time averages "+fmt.Sprintf("for epoch %v is %v", epoch, epochTimes2), 10)

	sort.Ints(epochTimes1)
	sequence := make([]int, 0)

	for i := 0; i < pr.numReplicas; i++ {
		found := false
		for j := 0; j < pr.numReplicas; j++ {
			if epochTimes1[i] == epochTimes2[j] {
				if pr.notInarray(sequence, j+1) {
					found = true
					sequence = append(sequence, j+1)
					break
				}
			}
		}
		if !found {
			panic("should not happen")
		}
	}

	pr.debug("leader ordering proposed for the epoch "+fmt.Sprintf("%v is %v", epoch, sequence), 10)

	s := ""

	for i := 0; i < pr.numReplicas; i++ {
		s = s + "," + strconv.FormatInt(int64(sequence[i]), 10)
	}

	return s[1:]
}

// checks if i is not in sequence

func (pr *Proxy) notInarray(sequence []int, i int) bool {
	for j := 0; j < len(sequence); j++ {
		if sequence[j] == i {
			return false
		}
	}
	return true
}
