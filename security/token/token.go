package token

import (
	"context"

	"github.com/alexfalkowski/go-service/crypto/argon2"
	"github.com/alexfalkowski/go-service/crypto/rand"
	"github.com/alexfalkowski/go-service/os"
)

// Generate key and hash for token.
func Generate() (string, string, error) {
	k, err := rand.GenerateString(32)
	if err != nil {
		return "", "", err
	}

	algo := argon2.NewAlgo()

	h, err := algo.Sign(k)
	if err != nil {
		return "", "", err
	}

	return k, h, nil
}

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
	}

	token struct {
		cfg  *Config
		algo argon2.Algo
	}

	none struct{}
)

// NewTokenizer for token.
func NewTokenizer(cfg *Config, algo argon2.Algo) Tokenizer {
	if !IsEnabled(cfg) {
		return &none{}
	}

	return &token{cfg: cfg, algo: algo}
}

// Generate token from secret file.
func (t *token) Generate(ctx context.Context) (context.Context, []byte, error) {
	d, err := os.ReadBase64File(t.cfg.Key)

	return ctx, []byte(d), err
}

// Verify the token with the stored hash.
func (t *token) Verify(ctx context.Context, token []byte) (context.Context, error) {
	return ctx, t.algo.Verify(t.cfg.Hash, string(token))
}

func (*none) Generate(ctx context.Context) (context.Context, []byte, error) {
	return ctx, nil, nil
}

func (*none) Verify(ctx context.Context, _ []byte) (context.Context, error) {
	return ctx, nil
}
