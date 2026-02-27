package pem

import (
	"encoding/pem"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/os"
)

var (
	// ErrInvalidBlock is returned when the provided data cannot be decoded into a PEM block.
	//
	// This typically indicates that the input is not PEM-encoded, is truncated, or does not contain
	// a recognizable PEM block header/footer.
	ErrInvalidBlock = errors.New("pem: invalid block")

	// ErrInvalidKind is returned when the decoded PEM block type does not match the expected kind.
	//
	// For example, if kind is "PUBLIC KEY" but the PEM block Type is "CERTIFICATE", Decode returns
	// ErrInvalidKind.
	ErrInvalidKind = errors.New("pem: invalid kind")
)

// NewDecoder constructs a Decoder that resolves and decodes PEM data using fs.
func NewDecoder(fs *os.FS) *Decoder {
	return &Decoder{fs}
}

// Decoder resolves PEM-encoded data and returns the raw bytes of a PEM block.
type Decoder struct {
	fs *os.FS
}

// Decode resolves PEM data from path (a go-service "source string") and returns the raw bytes for a PEM block of kind.
//
// The path parameter is resolved via os.FS.ReadSource, so it can be:
//   - "env:NAME" (read from environment variable NAME)
//   - "file:/path" (read from the filesystem)
//   - or a literal PEM value.
//
// Decode parses exactly the first PEM block in the resolved data. If the resolved value contains multiple
// PEM blocks, additional blocks are ignored.
//
// Errors:
//   - returns ErrInvalidBlock if no PEM block can be decoded.
//   - returns ErrInvalidKind if a PEM block is decoded but its Type does not equal kind.
//   - returns any error from fs.ReadSource if the input cannot be resolved/read.
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
