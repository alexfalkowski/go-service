package id

import "errors"

// ErrNotFound for id.
var ErrNotFound = errors.New("id: generator not found")

// Generator to generate an identifier.
type Generator interface {
	// Generate an identifier.
	Generate() string
}

// NewGenerator from config.
func NewGenerator(config *Config) (Generator, error) {
	if !IsEnabled(config) {
		return nil, nil
	}

	switch config.Kind {
	case "uuid":
		return &UUID{}, nil
	case "ksuid":
		return &KSUID{}, nil
	case "nanoid":
		return &NanoID{}, nil
	case "ulid":
		return &ULID{}, nil
	case "xid":
		return &XID{}, nil
	}

	return nil, ErrNotFound
}
