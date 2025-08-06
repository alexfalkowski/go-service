package hooks

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/encoding/base64"
	"github.com/alexfalkowski/go-service/v2/os"
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

// NewHook from config.
func NewHook(fs *os.FS, cfg *Config) (*hooks.Webhook, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	b, err := cfg.GetSecret(fs)
	if err != nil {
		return nil, err
	}

	return hooks.NewWebhook(bytes.String(b))
}
