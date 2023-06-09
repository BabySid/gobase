package gobase

import (
	"os"
	"time"
)

type FileWatcher interface {
	ChangeEvents(curFile string, size int64, newFiles []string) *FileEvents[EventMeta]
}

type EventMeta struct {
	NxtFile string
}

var _ FileWatcher = (*PollingFileWatcher)(nil)

type PollingFileWatcher struct {
}

const (
	pollInterval = 200 * time.Millisecond
)

func NewPollingFileWatcher() *PollingFileWatcher {
	fw := &PollingFileWatcher{}
	return fw
}

func (p *PollingFileWatcher) ChangeEvents(curFile string, size int64, newFiles []string) *FileEvents[EventMeta] {
	events := NewFileEvents[EventMeta]()

	go func() {
		prevSize := size

		curExist := true
		info, err := os.Stat(curFile)
		if err != nil {
			curExist = false
		}

		for {
			time.Sleep(pollInterval)

			newFile := ""
			for _, f := range newFiles {
				_, err := os.Stat(f)
				if err == nil {
					newFile = f
					break
				}
			}

			latest, err := os.Stat(curFile)
			if err != nil {
				if os.IsNotExist(err) {
					if newFile != "" {
						events.NotifyDeleted(EventMeta{NxtFile: newFile})
						return
					}
					continue
				}
			}

			if curExist && !os.SameFile(info, latest) {
				if newFile != "" {
					events.NotifyDeleted(EventMeta{NxtFile: newFile})
					return
				}
				continue
			}

			if prevSize > 0 && prevSize > latest.Size() {
				events.NotifyTruncated(EventMeta{NxtFile: curFile})
				return
			}

			if prevSize > 0 && prevSize < latest.Size() {
				events.NotifyModified(EventMeta{NxtFile: curFile})
				return
			}

			prevSize = latest.Size()
		}
	}()

	return events
}
