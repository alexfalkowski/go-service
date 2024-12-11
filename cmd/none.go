package cmd

import (
	"errors"
	"io/fs"
)

// ErrInvalidLocation for cmd.
var ErrInvalidLocation = errors.New("config: invalid location: (format kind:location)")

// None for cmd.
type None struct{}

// NewNone for cmd.
func NewNone() *None {
	return &None{}
}

// Read for none.
func (*None) Read() ([]byte, error) {
	return nil, ErrInvalidLocation
}

// Write for none.
func (*None) Write(_ []byte, _ fs.FileMode) error {
	return ErrInvalidLocation
}

// Kind for none.
func (*None) Kind() string {
	return "none"
}
