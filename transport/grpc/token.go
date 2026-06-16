package grpc

import (
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/token"
)

// NewToken constructs a token service for gRPC transport integration.
//
// The returned service is responsible for generating and verifying tokens according to cfg.Token (for example,
// JWT/PASETO/SSH token kinds as configured by the underlying token package).
//
// If cfg is disabled or cfg.Token is omitted, it returns nil so downstream wiring can treat token auth
// as not configured.
func NewToken(name env.Name, cfg *Config, fs *os.FS, gen id.Generator) *token.Token {
	if !cfg.IsEnabled() || !cfg.Token.IsEnabled() {
		return nil
	}
	return token.NewToken(name, cfg.Token, fs, gen)
}
