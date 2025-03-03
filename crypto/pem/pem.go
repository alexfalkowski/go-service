package pem

import (
	"encoding/pem"
	"errors"

	"github.com/alexfalkowski/go-service/os"
)

var (
	// ErrInvalidBlock of PEM.
	ErrInvalidBlock = errors.New("pem: invalid block")

	// ErrInvalidKind of PEM block.
	ErrInvalidKind = errors.New("pem: invalid kind")
)

// Decode from path.
func Decode(path, kind string) ([]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(b)
	if block == nil {
		return nil, ErrInvalidBlock
	}

	if block.Type != kind {
		return nil, ErrInvalidKind
	}

	return block.Bytes, nil
}
