package test

import (
	"github.com/alexfalkowski/go-sync"
)

// NewShutdowner returns a Shutdowner that records whether shutdown was requested.
func NewShutdowner() *Shutdowner {
	return &Shutdowner{ch: make(chan struct{})}
}

// Shutdowner is a test double that exposes application shutdown state to callers.
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

// Shutdown records a shutdown request and closes Done on first call.
func (s *Shutdowner) Shutdown(int) {
	s.m.Lock()
	defer s.m.Unlock()

	s.called = true
	s.once.Do(func() {
		close(s.ch)
	})
}
