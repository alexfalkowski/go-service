package argon2

import (
	"context"

	"github.com/alexfalkowski/go-service/crypto/argon2"
	"github.com/alexfalkowski/go-service/crypto/rand"
	"github.com/alexfalkowski/go-service/os"
)

// Tokenizer for argon2.
type Tokenizer struct {
	cfg  *Config
	algo argon2.Algo
}

// NewTokenizer for token.
func NewTokenizer(cfg *Config, algo argon2.Algo) *Tokenizer {
	return &Tokenizer{cfg: cfg, algo: algo}
}

// GenerateConfig for argon2.
func (t *Tokenizer) GenerateConfig() (string, string, error) {
	k, err := rand.GenerateString(32)
	if err != nil {
		return "", "", err
	}

	h, err := t.algo.Sign(k)
	if err != nil {
		return "", "", err
	}

	return k, h, nil
}

// Generate token from secret file.
func (t *Tokenizer) Generate(ctx context.Context) (context.Context, []byte, error) {
	d, err := os.ReadBase64File(t.cfg.Key)

	return ctx, []byte(d), err
}

// Verify the token with the stored hash.
func (t *Tokenizer) Verify(ctx context.Context, token []byte) (context.Context, error) {
	return ctx, t.algo.Verify(t.cfg.Hash, string(token))
}
