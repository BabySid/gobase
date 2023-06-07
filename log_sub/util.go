package log_sub

import (
	"github.com/BabySid/gobase"
	"time"
)

func verifyLogStep(layout string, name string) int {
	startTime, err := time.Parse(layout, name)
	gobase.True(err == nil)

	startTime = startTime.Add(time.Hour)
	if startTime.Format(layout) == name {
		return daily
	}
	return hourly
}
