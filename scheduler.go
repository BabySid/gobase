package gobase

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"sync"
)

type Scheduler struct {
	c    *cron.Cron
	lock sync.Mutex
	// name -> id
	jobs map[string]cron.EntryID
}

var GlobalScheduler *Scheduler

type ScheJob interface {
	Run()
}

func NewScheduler() *Scheduler {
	if GlobalScheduler != nil {
		return GlobalScheduler
	}

	GlobalScheduler = &Scheduler{
		c:    cron.New(cron.WithSeconds()),
		jobs: make(map[string]cron.EntryID),
	}
	return GlobalScheduler
}

func (s *Scheduler) AddJob(name string, spec string, cmd ScheJob) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.jobs[name]; ok {
		return fmt.Errorf("%s exist", name)
	}

	id, e := s.c.AddJob(spec, cmd)
	if e != nil {
		return e
	}

	s.jobs[name] = id
	return nil
}

func (s *Scheduler) DelJob(name string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	var id cron.EntryID
	id, ok := s.jobs[name]
	if !ok {
		return
	}

	s.c.Remove(id)

	delete(s.jobs, name)
}

func (s *Scheduler) Start() {
	s.c.Start()
}

func (s *Scheduler) Stop() {
	s.c.Stop()
}
