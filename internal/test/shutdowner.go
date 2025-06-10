package test

import (
	"sync"

	"github.com/alexfalkowski/go-service/v2/di"
)

func NewShutdowner() *Shutdowner {
	return &Shutdowner{}
}

type Shutdowner struct {
	called bool
	m      sync.RWMutex
}

func (s *Shutdowner) Called() bool {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.called
}

func (s *Shutdowner) Shutdown(...di.ShutdownOption) error {
	s.m.Lock()
	defer s.m.Unlock()

	s.called = true

	return nil
}
