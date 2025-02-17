package token

import (
	"bytes"
	"context"
	"fmt"
	"hash/crc32"
	"strconv"
	"strings"

	"github.com/alexfalkowski/go-service/crypto/rand"
	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/time"
)

const underscore = "_"

// Generate a token.
// The format is name_rand(64)_crc32(id).
func Generate(name env.Name, gen *rand.Generator) string {
	token := gen.Text()
	checksum := strconv.FormatUint(uint64(crc32.ChecksumIEEE([]byte(token))), 10)

	var builder strings.Builder

	builder.WriteString(string(name))
	builder.WriteString(underscore)
	builder.WriteString(token)
	builder.WriteString(underscore)
	builder.WriteString(checksum)

	return builder.String()
}

// Verify if the token matches the segments.
func Verify(name env.Name, token string) error {
	segments := strings.Split(token, underscore)

	if len(segments) != 3 {
		return fmt.Errorf("invalid length: %w", ErrInvalidMatch)
	}

	if segments[0] != string(name) {
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

// NewToken based on config.
func NewToken(cfg *Config, name env.Name, jwt *JWT, pas *Paseto) *Token {
	return &Token{cfg: cfg, name: name, jwt: jwt, pas: pas}
}

// Token will generate and verify based on what is defined in the config.
type Token struct {
	cfg  *Config
	jwt  *JWT
	pas  *Paseto
	name env.Name
}

func (t *Token) Generate(ctx context.Context) (context.Context, []byte, error) {
	if t.cfg == nil {
		return ctx, nil, nil
	}

	switch {
	case t.cfg.IsKey():
		d, err := os.ReadBase64File(t.cfg.Secret)

		return ctx, []byte(d), err
	case t.cfg.IsToken():
		d, err := os.ReadFile(t.cfg.Secret)

		return ctx, []byte(d), err
	case t.cfg.IsJWT():
		token, err := t.jwt.Generate(t.cfg.Subject, t.cfg.Audience, t.cfg.Issuer, time.MustParseDuration(t.cfg.Expiration))

		return ctx, []byte(token), err
	case t.cfg.IsPaseto():
		token, err := t.pas.Generate(t.cfg.Subject, t.cfg.Audience, t.cfg.Issuer, time.MustParseDuration(t.cfg.Expiration))

		return ctx, []byte(token), err
	}

	return ctx, nil, nil
}

func (t *Token) Verify(ctx context.Context, token []byte) (context.Context, error) {
	if t.cfg == nil {
		return ctx, nil
	}

	switch {
	case t.cfg.IsKey():
		d, err := os.ReadBase64File(t.cfg.Secret)
		if err != nil {
			return ctx, err
		}

		if !bytes.Equal([]byte(d), token) {
			return ctx, ErrInvalidMatch
		}

		return ctx, nil
	case t.cfg.IsToken():
		d, err := os.ReadFile(t.cfg.Secret)
		if err != nil {
			return ctx, err
		}

		if !bytes.Equal([]byte(d), token) {
			return ctx, ErrInvalidMatch
		}

		return ctx, Verify(t.name, string(token))
	case t.cfg.IsJWT():
		_, err := t.jwt.Verify(string(token), t.cfg.Audience, t.cfg.Issuer)

		return ctx, err
	case t.cfg.IsPaseto():
		_, err := t.pas.Verify(string(token), t.cfg.Audience, t.cfg.Issuer)

		return ctx, err
	}

	return ctx, nil
}
