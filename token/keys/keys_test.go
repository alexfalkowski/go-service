package keys_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/pem"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/token/errors"
	"github.com/alexfalkowski/go-service/v2/token/keys"
	"github.com/stretchr/testify/require"
)

func TestMapGet(t *testing.T) {
	key := &keys.Config{Config: test.NewEd25519()}

	tests := []struct {
		config keys.Map
		want   *keys.Config
		name   string
		id     string
	}{
		{name: "nil map", config: nil, id: "test"},
		{name: "missing key", config: keys.Map{"other": key}, id: "test"},
		{name: "nil key", config: keys.Map{"test": nil}, id: "test"},
		{name: "valid key", config: keys.Map{"test": key}, id: "test", want: key},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Same(t, tt.want, tt.config.Get(tt.id))
		})
	}
}

func TestConfigLoaders(t *testing.T) {
	loaders := []struct {
		load func(*keys.Config, *pem.Decoder) (any, error)
		name string
	}{
		{name: "signer", load: func(cfg *keys.Config, decoder *pem.Decoder) (any, error) {
			return cfg.Signer(decoder)
		}},
		{name: "verifier", load: func(cfg *keys.Config, decoder *pem.Decoder) (any, error) {
			return cfg.Verifier(decoder)
		}},
	}

	invalid := []struct {
		config  *keys.Config
		decoder *pem.Decoder
		name    string
	}{
		{name: "nil config", decoder: test.PEM},
		{name: "missing key config", config: &keys.Config{}, decoder: test.PEM},
		{name: "missing decoder", config: &keys.Config{Config: test.NewEd25519()}},
	}

	for _, loader := range loaders {
		t.Run(loader.name, func(t *testing.T) {
			for _, tt := range invalid {
				t.Run(tt.name, func(t *testing.T) {
					key, err := loader.load(tt.config, tt.decoder)
					require.Nil(t, key)
					require.ErrorIs(t, err, errors.ErrInvalidConfig)
				})
			}

			key, err := loader.load(&keys.Config{Config: test.NewEd25519()}, test.PEM)
			require.NoError(t, err)
			require.NotNil(t, key)
		})
	}
}
