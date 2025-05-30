package pem

import (
	"encoding/pem"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/os"
)

var (
	// ErrInvalidBlock of pem.
	ErrInvalidBlock = errors.New("pem: invalid block")

	// ErrInvalidKind of pem block.
	ErrInvalidKind = errors.New("pem: invalid kind")
)

// NewDecoder for pem.
func NewDecoder(fs *os.FS) *Decoder {
	return &Decoder{fs}
}

// Decoder for pem.
type Decoder struct {
	fs *os.FS
}

// Decode from path.
func (d *Decoder) Decode(path, kind string) ([]byte, error) {
	data, err := d.fs.ReadSource(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, ErrInvalidBlock
	}

	if block.Type != kind {
		return nil, ErrInvalidKind
	}

	return block.Bytes, nil
}
