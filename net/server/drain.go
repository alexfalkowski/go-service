package server

import (
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-sync"
)

// ErrDraining is returned when server shutdown has started.
var ErrDraining = errors.New("server: draining")

// Drain tracks whether server shutdown has started.
type Drain struct {
	done     chan struct{}
	draining sync.Bool
}

// NewDrain creates a drain state tracker.
func NewDrain() *Drain {
	return &Drain{done: make(chan struct{})}
}

// Start marks the server lifecycle as draining.
//
// Start must be called once for a drain instance.
func (d *Drain) Start() {
	d.draining.Store(true)
	close(d.done)
}

// Error returns ErrDraining after Start has been called.
func (d *Drain) Error() error {
	if !d.draining.Load() {
		return nil
	}

	return ErrDraining
}

// Done returns a channel that is closed after Start has been called.
func (d *Drain) Done() <-chan struct{} {
	return d.done
}
