package test

import (
	"sync"

	"go.uber.org/fx"
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

func (s *Shutdowner) Shutdown(...fx.ShutdownOption) error {
	s.m.Lock()
	defer s.m.Unlock()

	s.called = true

	return nil
}
