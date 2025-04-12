package io

import "github.com/alexfalkowski/go-service/os"

// NewNone for io.
func NewNone() *None {
	return &None{}
}

// None will just fail if you try to read or write to it.
type None struct{}

// Read fails with invalid location.
func (*None) Read() ([]byte, error) {
	return nil, ErrInvalidLocation
}

// Write fails with invalid location.
func (*None) Write(_ []byte, _ os.FileMode) error {
	return ErrInvalidLocation
}

// Kind is invalid.
func (*None) Kind() string {
	return "none"
}
