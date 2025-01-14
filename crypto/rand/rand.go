package rand

import (
	"crypto/rand"
	"io"
	"math/big"
)

type (
	// Reader is just rand.Reader.
	Reader io.Reader

	// Generator for rand.
	Generator struct {
		reader  Reader
		letters string
		symbols string
	}
)

// NewReader for rand.
func NewReader() Reader {
	return rand.Reader
}

// NewGenerator for rand.
func NewGenerator(reader Reader) *Generator {
	return &Generator{
		letters: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
		symbols: "~!@#$%^&*()_+-={}|[]<>?,./",
		reader:  reader,
	}
}

// Read for rand.
func (g *Generator) Read(b []byte) (int, error) {
	return io.ReadFull(g.reader, b)
}

// GenerateBytes for rand.
func (g *Generator) GenerateBytes(size uint32) ([]byte, error) {
	bytes := make([]byte, size)
	_, err := g.Read(bytes)

	return bytes, err
}

// GenerateString will generate using letters and symbols.
func (g *Generator) GenerateString(size uint32) (string, error) {
	return g.generate(size, g.letters+g.symbols)
}

// GenerateLetters will generate using letters.
func (g *Generator) GenerateLetters(size uint32) (string, error) {
	return g.generate(size, g.letters)
}

func (g *Generator) generate(size uint32, values string) (string, error) {
	bytes := make([]byte, size)

	for i := range size {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(values))))
		if err != nil {
			return "", err
		}

		bytes[i] = values[num.Int64()]
	}

	return string(bytes), nil
}
