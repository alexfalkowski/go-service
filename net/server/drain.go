package server

import (
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-sync"
)

// ErrDraining is returned when server shutdown has started.
var ErrDraining = errors.New("server: draining")

// Drain tracks whether server shutdown has started.
type Drain struct {
	draining sync.Bool
}

// NewDrain creates a drain state tracker.
func NewDrain() *Drain {
	return &Drain{}
}

// Start marks the server lifecycle as draining.
func (d *Drain) Start() {
	if d == nil {
		return
	}

	d.draining.Store(true)
}

// Error returns ErrDraining after Start has been called.
func (d *Drain) Error() error {
	if d == nil || !d.draining.Load() {
		return nil
	}

	return ErrDraining
}
