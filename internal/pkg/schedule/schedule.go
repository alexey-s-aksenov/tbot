package schedule

import (
	"time"
)

//FirstStart calculates time to next hour
func FirstStart(hour int) time.Duration {
	const (
		unit = time.Hour
		//layout = "15:04:05.000"
	)
	var currentTime = time.Now()

	var startTime = currentTime.Truncate(unit)
	currentHour := startTime.Hour()
	var addHours int
	if currentHour < hour {
		addHours = hour - currentHour
	} else {
		addHours = 24 - (currentHour - hour)
	}
	startTime = startTime.Add(time.Duration(addHours) * unit)
	return startTime.Sub(currentTime)
}
