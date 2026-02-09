package test

import (
	"sync"

	"github.com/alexfalkowski/go-service/v2/di"
)

func NewShutdowner() *Shutdowner {
	return &Shutdowner{ch: make(chan struct{})}
}

type Shutdowner struct {
	ch     chan struct{}
	m      sync.RWMutex
	once   sync.Once
	called bool
}

func (s *Shutdowner) Called() bool {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.called
}

func (s *Shutdowner) Done() <-chan struct{} {
	return s.ch
}

func (s *Shutdowner) Shutdown(...di.ShutdownOption) error {
	s.m.Lock()
	defer s.m.Unlock()

	s.called = true
	s.once.Do(func() {
		close(s.ch)
	})

	return nil
}
