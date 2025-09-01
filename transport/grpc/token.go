package grpc

import (
	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/token"
)

// .NewController for gRPC.
func NewController(cfg *Config) (token.AccessController, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}
	return token.NewAccessController(cfg.Token)
}

// NewToken for gRPC.
func NewToken(name env.Name, cfg *Config, fs *os.FS, sig *ed25519.Signer, ver *ed25519.Verifier, gen id.Generator) *token.Token {
	if !cfg.IsEnabled() {
		return nil
	}
	return token.NewToken(name, cfg.Token, fs, sig, ver, gen)
}
