package hooks

import (
	"github.com/alexfalkowski/go-service/bytes"
	"github.com/alexfalkowski/go-service/crypto/rand"
	"github.com/alexfalkowski/go-service/encoding/base64"
	"github.com/alexfalkowski/go-service/os"
	hooks "github.com/standard-webhooks/standard-webhooks/libraries/go"
)

// NewGenerator for hooks.
func NewGenerator(gen *rand.Generator) *Generator {
	return &Generator{gen: gen}
}

// Generator for hooks.
type Generator struct {
	gen *rand.Generator
}

// Generate for hooks.
func (g *Generator) Generate() (string, error) {
	b, err := g.gen.GenerateBytes(32)

	return base64.Encode(b), err
}

// New hook from config.
func New(fs *os.FS, cfg *Config) (*hooks.Webhook, error) {
	if !IsEnabled(cfg) {
		return nil, nil
	}

	b, err := cfg.GetSecret(fs)
	if err != nil {
		return nil, err
	}

	return hooks.NewWebhook(bytes.String(b))
}
