package opaque

import (
	"fmt"
	"hash/crc32"
	"strconv"
	"strings"

	"github.com/alexfalkowski/go-service/crypto/rand"
	"github.com/alexfalkowski/go-service/token/errors"
)

const underscore = "_"

// NewToken creates a new opaque token.
func NewToken(cfg *Config, generator *rand.Generator) *Token {
	if !IsEnabled(cfg) {
		return nil
	}

	return &Token{cfg: cfg, generator: generator}
}

// Token represents an opaque token.
type Token struct {
	cfg       *Config
	generator *rand.Generator
}

// Generate generates a new opaque token.
func (t *Token) Generate() string {
	token := t.generator.Text()
	checksum := strconv.FormatUint(uint64(crc32.ChecksumIEEE([]byte(token))), 10)

	var builder strings.Builder

	builder.WriteString(t.cfg.Name)
	builder.WriteString(underscore)
	builder.WriteString(token)
	builder.WriteString(underscore)
	builder.WriteString(checksum)

	return builder.String()
}

// Verify verifies an opaque token.
func (t *Token) Verify(token string) error {
	segments := strings.Split(token, underscore)

	if len(segments) != 3 {
		return fmt.Errorf("invalid length: %w", errors.ErrInvalidMatch)
	}

	if segments[0] != t.cfg.Name {
		return fmt.Errorf("invalid prefix: %w", errors.ErrInvalidMatch)
	}

	u64, err := strconv.ParseUint(segments[2], 10, 32)
	if err != nil {
		return fmt.Errorf("%w: %w", err, errors.ErrInvalidMatch)
	}

	if crc32.ChecksumIEEE([]byte(segments[1])) != uint32(u64) {
		return fmt.Errorf("invalid checksum: %w", errors.ErrInvalidMatch)
	}

	return nil
}
