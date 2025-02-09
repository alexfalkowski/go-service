package cmd

import (
	"errors"

	"github.com/alexfalkowski/go-service/os"
)

// ErrInvalidLocation for cmd.
var ErrInvalidLocation = errors.New("invalid location (format kind:location)")

// None for cmd.
type None struct{}

// NewNone for cmd.
func NewNone() *None {
	return &None{}
}

// Read for none.
func (*None) Read() (string, error) {
	return "", ErrInvalidLocation
}

// Write for none.
func (*None) Write(_ string, _ os.FileMode) error {
	return ErrInvalidLocation
}

// Kind for none.
func (*None) Kind() string {
	return "none"
}
