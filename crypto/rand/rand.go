package rand

import (
	"crypto/rand"
	"io"

	"github.com/samber/lo"
)

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

// Text returns a cryptographically random string.
func (g *Generator) Text() string {
	return rand.Text()
}

// GenerateBytes returns a cryptographically random byte slice of size.
func (g *Generator) GenerateBytes(size int) ([]byte, error) {
	bytes := make([]byte, size)
	_, err := g.Read(bytes)

	return bytes, err
}

// GenerateText will generate using alphanumeric charset.
func (g *Generator) GenerateText(size int) string {
	return lo.RandomString(size, lo.AlphanumericCharset)
}
