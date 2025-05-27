package io

import "github.com/alexfalkowski/go-service/v2/errors"

var (
	// ErrLocationMissing for cmd.
	ErrLocationMissing = errors.New("location is missing")

	// ErrInvalidLocation for io.
	ErrInvalidLocation = errors.New("invalid location (format kind:location)")
)
