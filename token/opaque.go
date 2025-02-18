package token

import (
	"fmt"
	"hash/crc32"
	"strconv"
	"strings"

	"github.com/alexfalkowski/go-service/crypto/rand"
	"github.com/alexfalkowski/go-service/env"
)

const underscore = "_"

// NewOpaque creates a new opaque token.
func NewOpaque(name env.Name, generator *rand.Generator) *Opaque {
	return &Opaque{name: name, generator: generator}
}

// Opaque represents an opaque token.
type Opaque struct {
	generator *rand.Generator
	name      env.Name
}

// Generate generates a new opaque token.
func (o *Opaque) Generate() string {
	token := o.generator.Text()
	checksum := strconv.FormatUint(uint64(crc32.ChecksumIEEE([]byte(token))), 10)

	var builder strings.Builder

	builder.WriteString(o.name.String())
	builder.WriteString(underscore)
	builder.WriteString(token)
	builder.WriteString(underscore)
	builder.WriteString(checksum)

	return builder.String()
}

// Verify verifies an opaque token.
func (o *Opaque) Verify(token string) error {
	segments := strings.Split(token, underscore)

	if len(segments) != 3 {
		return fmt.Errorf("invalid length: %w", ErrInvalidMatch)
	}

	if segments[0] != o.name.String() {
		return fmt.Errorf("invalid prefix: %w", ErrInvalidMatch)
	}

	u64, err := strconv.ParseUint(segments[2], 10, 32)
	if err != nil {
		return fmt.Errorf("%w: %w", err, ErrInvalidMatch)
	}

	if crc32.ChecksumIEEE([]byte(segments[1])) != uint32(u64) {
		return fmt.Errorf("invalid checksum: %w", ErrInvalidMatch)
	}

	return nil
}
