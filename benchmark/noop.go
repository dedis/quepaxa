package benchmark

import "time"

type NoOpApplication struct {
	SleepDuration int64 // sleep time in micro seconds
}

func (sleepApp *NoOpApplication) executeNoOpApp(request string) string {
	start := time.Now()
	for time.Now().Sub(start).Nanoseconds() < sleepApp.SleepDuration*1000 {
		// busy waiting
	}
	return request
}
