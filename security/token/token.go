package token

import (
	"context"
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

	// NoopTokenizer for token.
	NoopTokenizer struct{}
)

func (*NoopTokenizer) GenerateConfig() (string, string, error) {
	return "", "", nil
}

func (*NoopTokenizer) Generate(ctx context.Context) (context.Context, []byte, error) {
	return ctx, nil, nil
}

func (*NoopTokenizer) Verify(ctx context.Context, _ []byte) (context.Context, error) {
	return ctx, nil
}
