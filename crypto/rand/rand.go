package rand

import (
	"crypto/rand"
	"io"
	"math/big"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/strings"
)

const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// NewReader returns a cryptographically secure random Reader.
//
// This is a thin wrapper around crypto/rand.Reader.
func NewReader() Reader {
	return rand.Reader
}

// Reader is a cryptographically secure random source.
//
// It is defined as an aliasable type so it can be injected/mocked in tests while still behaving like
// an io.Reader.
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
	data := make([]byte, size)
	_, err := g.Read(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GenerateText returns a cryptographically random string of length size.
//
// Characters are drawn from the package's letter set, making this helper
// suitable for text tokens but not a substitute for GenerateBytes when binary
// randomness is required.
func (g *Generator) GenerateText(size int) (string, error) {
	data := make([]byte, size)
	length := int64(len(letters))

	for i := range size {
		num, err := rand.Int(g.reader, big.NewInt(length))
		if err != nil {
			return strings.Empty, err
		}

		data[i] = letters[num.Int64()]
	}

	return bytes.String(data), nil
}
