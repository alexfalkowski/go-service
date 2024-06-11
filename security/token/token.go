package token

import (
	"context"

	"github.com/alexfalkowski/go-service/crypto/argon2"
	ta "github.com/alexfalkowski/go-service/security/token/argon2"
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

	// Tokenizer will generate and verify.
	Tokenizer interface {
		Generator
		Verifier

		// GenerateConfig to be used by the Tokenizer.
		GenerateConfig() (string, string, error)
	}

	none struct{}
)

// NewTokenizer for token.
func NewTokenizer(cfg *Config, algo argon2.Algo) Tokenizer {
	if !IsEnabled(cfg) {
		return &none{}
	}

	return ta.NewTokenizer(cfg.Argon2, algo)
}

func (*none) GenerateConfig() (string, string, error) {
	return "", "", nil
}

func (*none) Generate(ctx context.Context) (context.Context, []byte, error) {
	return ctx, nil, nil
}

func (*none) Verify(ctx context.Context, _ []byte) (context.Context, error) {
	return ctx, nil
}
