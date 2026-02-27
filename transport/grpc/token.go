package grpc

import (
	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/token"
)

// NewController constructs an access controller for gRPC token authorization.
//
// When token auth is enabled (via cfg.Token), the returned controller can be used to authorize
// authenticated subjects against configured access rules.
//
// If cfg is disabled, it returns (nil, nil) so downstream wiring can treat access control as not configured.
func NewController(cfg *Config) (token.AccessController, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}
	return token.NewAccessController(cfg.Token)
}

// NewToken constructs a token service for gRPC transport integration.
//
// The returned service is responsible for generating and verifying tokens according to cfg.Token (for example,
// JWT/PASETO/SSH token kinds as configured by the underlying token package).
//
// If cfg is disabled, it returns nil so downstream wiring can treat token auth as not configured.
func NewToken(name env.Name, cfg *Config, fs *os.FS, sig *ed25519.Signer, ver *ed25519.Verifier, gen id.Generator) *token.Token {
	if !cfg.IsEnabled() {
		return nil
	}
	return token.NewToken(name, cfg.Token, fs, sig, ver, gen)
}
