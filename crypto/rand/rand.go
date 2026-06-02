package rand

import (
	"crypto/fips140"
	"crypto/rand"
	"math/big"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/strings"
)

const alphanumeric = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// ErrInvalidSize is returned when a random value size is negative.
var ErrInvalidSize = errors.New("rand: invalid size")

// NewReader returns a cryptographically secure random Reader.
//
// This is a thin wrapper around crypto/rand.Reader.
func NewReader() Reader {
	if fips140.Enforced() {
		return rand.Reader
	}

	return Reader(rand.Reader)
}

// Reader is a cryptographically secure random source.
//
// It is intentionally the same shape as io.Reader while remaining part of this package's crypto API.
type Reader io.Reader

// NewGenerator constructs a Generator that draws randomness from reader.
//
// The provided reader should be cryptographically secure (for example the value returned by NewReader).
func NewGenerator(reader Reader) *Generator {
	return &Generator{reader: reader}
}

// Generator produces cryptographically secure random values.
//
// It provides both raw byte reads and convenience helpers for generating random
// bytes or random text.
type Generator struct {
	reader Reader
}

// Reader returns the source used by the Generator.
func (g *Generator) Reader() Reader {
	return g.reader
}

// Read fills b with random bytes read from the underlying Reader.
//
// It reads exactly len(b) bytes or returns an error.
func (g *Generator) Read(b []byte) (int, error) {
	return io.ReadFull(g.reader, b)
}

// GenerateBytes returns a cryptographically random byte slice of length size.
//
// The returned bytes are read directly from the underlying Reader and may span
// the full 0-255 byte range.
func (g *Generator) GenerateBytes(size int) ([]byte, error) {
	if size < 0 {
		return nil, ErrInvalidSize
	}

	data := make([]byte, size)
	_, err := g.Read(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GenerateText returns a cryptographically random string of length size.
//
// Characters are drawn from the package's alphanumeric set, making this helper
// suitable for text tokens but not a substitute for GenerateBytes when binary
// randomness is required.
func (g *Generator) GenerateText(size int) (string, error) {
	if size < 0 {
		return strings.Empty, ErrInvalidSize
	}

	data := make([]byte, size)
	length := int64(len(alphanumeric))

	for i := range size {
		num, err := rand.Int(g.reader, big.NewInt(length))
		if err != nil {
			return strings.Empty, err
		}

		data[i] = alphanumeric[num.Int64()]
	}

	return bytes.String(data), nil
}
