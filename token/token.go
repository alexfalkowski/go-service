package token

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"hash/crc32"
	"strconv"
	"strings"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/runtime"
	st "github.com/alexfalkowski/go-service/time"
	"github.com/google/uuid"
)

const underscore = "_"

// Generate a token.
// The format is os.ExecutableName_uuid.NewV7_crc32(uuid).
func Generate(name env.Name) string {
	id, err := uuid.NewV7()
	runtime.Must(err)

	uuid := id.String()
	checksum := strconv.FormatUint(uint64(crc32.ChecksumIEEE([]byte(uuid))), 10)

	var builder strings.Builder

	builder.WriteString(string(name))
	builder.WriteString(underscore)
	builder.WriteString(uuid)
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

	if _, err := uuid.Parse(segments[1]); err != nil {
		return fmt.Errorf("%w: %w", err, ErrInvalidMatch)
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

var (
	// ErrInvalidMatch for token.
	ErrInvalidMatch = errors.New("token: invalid match")

	// ErrInvalidIssuer for service.
	ErrInvalidIssuer = errors.New("token: invalid issuer")

	// ErrInvalidAudience for service.
	ErrInvalidAudience = errors.New("token: invalid audience")

	// ErrInvalidTime for service.
	ErrInvalidTime = errors.New("token: invalid time")
)

type (
	// Generator allows the implementation of different types generators.
	Generator interface {
		// Generate a new token or error.
		Generate(ctx context.Context) (context.Context, []byte, error)
	}

	// Verifier allows the implementation of different types of verifiers.
	Verifier interface {
		// Verify a token or error.
		Verify(ctx context.Context, token []byte) (context.Context, error)
	}

	// Token will generate and verify based on what is defined in the config.
	Token struct {
		cfg  *Config
		jwt  *JWT
		pas  *Paseto
		name env.Name
	}
)

// NewToken based on config.
func NewToken(cfg *Config, name env.Name, jwt *JWT, pas *Paseto) *Token {
	return &Token{cfg: cfg, name: name, jwt: jwt, pas: pas}
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
		token, err := t.jwt.Generate(t.cfg.Subject, t.cfg.Audience, t.cfg.Issuer, st.MustParseDuration(t.cfg.Expiration))

		return ctx, []byte(token), err
	case t.cfg.IsPaseto():
		token, err := t.pas.Generate(t.cfg.Subject, t.cfg.Audience, t.cfg.Issuer, st.MustParseDuration(t.cfg.Expiration))

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
