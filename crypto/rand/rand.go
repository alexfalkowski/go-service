package rand

import (
	"crypto/rand"
	"io"
	"math/big"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/strings"
)

const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// NewReader for rand.
func NewReader() Reader {
	return rand.Reader
}

// Reader is just rand.Reader.
type Reader io.Reader

// NewGenerator for rand.
func NewGenerator(reader Reader) *Generator {
	return &Generator{reader: reader}
}

// Generator for rand.
type Generator struct {
	reader Reader
}

// Read for rand.
func (g *Generator) Read(b []byte) (int, error) {
	return io.ReadFull(g.reader, b)
}

// GenerateBytes returns a cryptographically random byte slice of size.
func (g *Generator) GenerateBytes(size int) ([]byte, error) {
	s, err := g.generate(size, letters)
	return strings.Bytes(s), err
}

// GenerateText will generate using letters.
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
