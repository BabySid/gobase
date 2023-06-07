package gobase

import (
	"os"
	"time"
)

type FileWatcher interface {
	ChangeEvents(size int64) (*FileEvents, error)
}

var _ FileWatcher = (*PollingFileWatcher)(nil)

type PollingFileWatcher struct {
	curFile string
	newFile string
	Size    int64
}

const (
	pollInterval = 200 * time.Millisecond
)

func NewPollingFileWatcher(curFile string, newFile string) *PollingFileWatcher {
	fw := &PollingFileWatcher{curFile: curFile, newFile: newFile, Size: 0}
	return fw
}

func (p *PollingFileWatcher) ChangeEvents(size int64) (*FileEvents, error) {
	info, err := os.Stat(p.curFile)
	if err != nil {
		return nil, err
	}

	p.Size = size

	events := NewFileEvents()

	go func() {
		prevSize := p.Size
		for {
			time.Sleep(pollInterval)

			latest, err := os.Stat(p.curFile)
			if err != nil {
				if os.IsNotExist(err) {
					events.NotifyDeleted()
					return
				}
			}

			if !os.SameFile(info, latest) {
				events.NotifyDeleted()
				return
			}

			p.Size = latest.Size()
			if prevSize > 0 && prevSize > p.Size {
				events.NotifyTruncated()
				return
			}

			if prevSize > 0 && prevSize < p.Size {
				events.NotifyModified()
				return
			}

			prevSize = p.Size

			_, err = os.Stat(p.newFile)
			if err == nil {
				events.NotifyCreated()
				return
			}
		}
	}()

	return events, nil
}
