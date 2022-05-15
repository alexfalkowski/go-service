package test

import (
	"sync"

	"go.uber.org/fx"
)

func NewShutdowner() *Shutdowner {
	return &Shutdowner{}
}

type Shutdowner struct {
	mux    sync.RWMutex
	called bool
}

func (s *Shutdowner) Called() bool {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return s.called
}

func (s *Shutdowner) Shutdown(...fx.ShutdownOption) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.called = true

	return nil
}
