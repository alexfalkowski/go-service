package test

import (
	"sync"

	"github.com/alexfalkowski/go-service/v2/di"
)

// NewShutdowner creates a Shutdowner for tests.
func NewShutdowner() *Shutdowner {
	return &Shutdowner{ch: make(chan struct{})}
}

// Shutdowner for tests.
type Shutdowner struct {
	ch     chan struct{}
	m      sync.RWMutex
	once   sync.Once
	called bool
}

// Called reports whether Shutdown has been invoked.
func (s *Shutdowner) Called() bool {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.called
}

// Done returns a channel that is closed when Shutdown is invoked.
func (s *Shutdowner) Done() <-chan struct{} {
	return s.ch
}

// Shutdown implements di.Shutdowner and closes Done on first call.
func (s *Shutdowner) Shutdown(...di.ShutdownOption) error {
	s.m.Lock()
	defer s.m.Unlock()

	s.called = true
	s.once.Do(func() {
		close(s.ch)
	})

	return nil
}
