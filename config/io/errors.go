package io

import "errors"

var (
	// ErrLocationMissing for cmd.
	ErrLocationMissing = errors.New("location is missing")

	// ErrInvalidLocation for io.
	ErrInvalidLocation = errors.New("invalid location (format kind:location)")
)
