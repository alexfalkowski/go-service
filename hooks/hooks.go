package hooks

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/encoding/base64"
	"github.com/alexfalkowski/go-service/v2/os"
	hooks "github.com/standard-webhooks/standard-webhooks/libraries/go"
)

// NewGenerator constructs a Generator for creating webhook secrets.
func NewGenerator(gen *rand.Generator) *Generator {
	return &Generator{gen: gen}
}

// Generator generates secrets suitable for Standard Webhooks.
type Generator struct {
	gen *rand.Generator
}

// Generate returns a new base64-encoded secret.
//
// It generates 32 random bytes and encodes them as base64.
func (g *Generator) Generate() (string, error) {
	b, err := g.gen.GenerateBytes(32)
	return base64.Encode(b), err
}

// NewHook constructs a Standard Webhooks hook instance from cfg.
//
// If cfg is disabled, it returns (nil, nil). Otherwise it loads the configured secret (via Config.GetSecret)
// and uses it to construct the hook instance.
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
