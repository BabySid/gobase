package gobase

type FileEvents struct {
	Created   chan bool
	Modified  chan bool
	Truncated chan bool
	Deleted   chan bool // or renames
}

func NewFileEvents() *FileEvents {
	return &FileEvents{
		Created:   make(chan bool, 1),
		Modified:  make(chan bool, 1),
		Truncated: make(chan bool, 1),
		Deleted:   make(chan bool, 1),
	}
}

func (fe *FileEvents) NotifyCreated() {
	sendOnlyIfEmpty(fe.Created)
}

func (fe *FileEvents) NotifyModified() {
	sendOnlyIfEmpty(fe.Modified)
}

func (fe *FileEvents) NotifyTruncated() {
	sendOnlyIfEmpty(fe.Truncated)
}

func (fe *FileEvents) NotifyDeleted() {
	sendOnlyIfEmpty(fe.Deleted)
}

// sendOnlyIfEmpty sends on a bool channel only if the channel has no
// backlog to be read by other goroutines. This concurrency pattern
// can be used to notify other goroutines if and only if they are
// looking for it (i.e., subsequent notifications can be compressed
// into one).
func sendOnlyIfEmpty(ch chan bool) {
	select {
	case ch <- true:
	default:
	}
}
