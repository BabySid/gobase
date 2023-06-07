package gobase

type FileEvents struct {
	Modified  chan bool
	Truncated chan bool
	Deleted   chan bool // or renames
}

func NewFileEvents() *FileEvents {
	return &FileEvents{
		make(chan bool, 1), make(chan bool, 1), make(chan bool, 1)}
}

func (fc *FileEvents) NotifyModified() {
	sendOnlyIfEmpty(fc.Modified)
}

func (fc *FileEvents) NotifyTruncated() {
	sendOnlyIfEmpty(fc.Truncated)
}

func (fc *FileEvents) NotifyDeleted() {
	sendOnlyIfEmpty(fc.Deleted)
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
