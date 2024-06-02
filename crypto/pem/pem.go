package pem

import (
	"encoding/pem"
	"errors"
	"os"
)

var (
	// ErrInvalidBlock of PEM.
	ErrInvalidBlock = errors.New("invalid block")

	// ErrInvalidKind of PEM block.
	ErrInvalidKind = errors.New("invalid kind")
)

// Decode from path.
func Decode(path, kind string) ([]byte, error) {
	d, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(d)
	if block == nil {
		return nil, ErrInvalidBlock
	}

	if block.Type != kind {
		return nil, ErrInvalidKind
	}

	return block.Bytes, nil
}
