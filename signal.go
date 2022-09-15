package gobase

import (
	"os"
	"os/signal"
	"sync"
)

type SignalHandler func(sig os.Signal)

type SignalSet struct {
	handles sync.Map
}

func NewSignalSet() *SignalSet {
	s := &SignalSet{handles: sync.Map{}}
	go s.run()
	return s
}

func (s *SignalSet) Register(sig os.Signal, handler SignalHandler) {
	s.handles.Store(sig, handler)
}

func (s *SignalSet) run() {
	c := make(chan os.Signal, 1)
	signal.Notify(c)

	for sig := range c {
		if h, ok := s.handles.Load(sig); ok {
			h.(SignalHandler)(sig)
		}
	}
}
