package benchmark

import "time"

/*
No Op application implements an echo service where the client response is equal to the request. It also defines a synthetic service time, where the response is delayed until the time
*/

type NoOpApplication struct {
	SleepDuration int64 // sleep time in micro seconds
}

func (sleepApp *NoOpApplication) executeNoOpApp(request string) string {
	start := time.Now()
	for time.Now().Sub(start).Microseconds() < sleepApp.SleepDuration {
		// busy waiting for the sleep time to complete
	}
	return request
}
