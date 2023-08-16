package log_sub

import (
	"github.com/BabySid/gobase"
	"path/filepath"
	"time"
)

func verifyLogStep(layout DateTimeLayout, name string) int {
	fName := filepath.Base(name)
	startTime, err := time.Parse(layout.Layout, fName)
	gobase.True(err == nil, "parse(%s %s) failed. err=%v", layout.Layout, fName, err)

	startTime = startTime.Add(time.Hour)
	if startTime.Format(layout.Layout) == fName {
		return daily
	}
	return hourly
}
