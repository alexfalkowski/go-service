package pem

import (
	"encoding/pem"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/os"
)

var (
	// ErrInvalidBlock is returned when PEM decoding fails to find a PEM block.
	ErrInvalidBlock = errors.New("pem: invalid block")

	// ErrInvalidKind is returned when the decoded PEM block type does not match the expected kind.
	ErrInvalidKind = errors.New("pem: invalid kind")
)

// NewDecoder constructs a Decoder that reads PEM data using fs.
func NewDecoder(fs *os.FS) *Decoder {
	return &Decoder{fs}
}

// Decoder reads PEM-encoded data and returns the raw bytes of a PEM block.
type Decoder struct {
	fs *os.FS
}

// Decode reads PEM data from path (a go-service "source string") and returns the raw bytes for a PEM block of kind.
//
// path is read via os.FS.ReadSource, so it can be:
//   - "env:NAME" (read from environment variable NAME)
//   - "file:/path" (read from file system)
//   - or a literal PEM value.
//
// It returns ErrInvalidBlock if the data does not decode into a PEM block.
// It returns ErrInvalidKind if the decoded block type does not equal kind.
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
