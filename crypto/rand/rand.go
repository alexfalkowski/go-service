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
func NewReader() Reader {
	return rand.Reader
}

// Reader is a cryptographically secure random source.
type Reader io.Reader

// NewGenerator constructs a Generator that draws randomness from reader.
func NewGenerator(reader Reader) *Generator {
	return &Generator{reader: reader}
}

// Generator produces cryptographically secure random values.
type Generator struct {
	reader Reader
}

// Read fills b with random bytes read from the underlying Reader.
func (g *Generator) Read(b []byte) (int, error) {
	return io.ReadFull(g.reader, b)
}

// GenerateBytes returns a cryptographically random byte slice of length size.
//
// The bytes are drawn by generating random characters from the package's letter set and converting the
// resulting string to bytes. The returned slice must be treated as read-only.
func (g *Generator) GenerateBytes(size int) ([]byte, error) {
	s, err := g.generate(size, letters)
	return strings.Bytes(s), err
}

// GenerateText returns a cryptographically random string of length size drawn from the package's letter set.
func (g *Generator) GenerateText(size int) (string, error) {
	return g.generate(size, letters)
}

func (g *Generator) generate(size int, values string) (string, error) {
	data := make([]byte, size)
	length := int64(len(values))

	for i := range size {
		num, err := rand.Int(g.reader, big.NewInt(length))
		if err != nil {
			return strings.Empty, err
		}

		data[i] = values[num.Int64()]
	}

	return bytes.String(data), nil
}
