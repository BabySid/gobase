package gobase

type FileEvents[T any] struct {
	Created   chan T
	Modified  chan T
	Truncated chan T
	Deleted   chan T // or renames
}

func NewFileEvents[T any]() *FileEvents[T] {
	return &FileEvents[T]{
		Created:   make(chan T, 1),
		Modified:  make(chan T, 1),
		Truncated: make(chan T, 1),
		Deleted:   make(chan T, 1),
	}
}

func (fe *FileEvents[T]) NotifyCreated(meta T) {
	sendOnlyIfEmpty(fe.Created, meta)
}

func (fe *FileEvents[T]) NotifyModified(meta T) {
	sendOnlyIfEmpty(fe.Modified, meta)
}

func (fe *FileEvents[T]) NotifyTruncated(meta T) {
	sendOnlyIfEmpty(fe.Truncated, meta)
}

func (fe *FileEvents[T]) NotifyDeleted(meta T) {
	sendOnlyIfEmpty(fe.Deleted, meta)
}

// sendOnlyIfEmpty sends on a bool channel only if the channel has no
// backlog to be read by other goroutines. This concurrency pattern
// can be used to notify other goroutines if and only if they are
// looking for it (i.e., subsequent notifications can be compressed
// into one).
func sendOnlyIfEmpty[T any](ch chan T, meta T) {
	select {
	case ch <- meta:
	default:
	}
}
