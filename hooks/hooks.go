package hooks

import (
	"encoding/base64"

	"github.com/alexfalkowski/go-service/crypto/rand"
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

	return base64.StdEncoding.EncodeToString(b), err
}

// New hook from config.
func New(cfg *Config) (*hooks.Webhook, error) {
	if cfg == nil {
		return hooks.NewWebhookRaw(nil)
	}

	b, err := cfg.GetSecret()
	if err != nil {
		return nil, err
	}

	return hooks.NewWebhook(string(b))
}
