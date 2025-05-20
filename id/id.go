package id

import (
	"errors"

	"github.com/alexfalkowski/go-service/v2/crypto/rand"
)

// ErrNotFound for id.
var ErrNotFound = errors.New("id: generator not found")

// Generator to generate an identifier.
type Generator interface {
	// Generate an identifier.
	Generate() string
}

// NewGenerator from config.
func NewGenerator(config *Config, reader rand.Reader) (Generator, error) {
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
		return NewULID(reader), nil
	case "xid":
		return &XID{}, nil
	}

	return nil, ErrNotFound
}
