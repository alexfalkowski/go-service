package feature

import (
	"github.com/alexfalkowski/go-service/os"
	flipt "github.com/open-feature/go-sdk-contrib/providers/flipt/pkg/provider"
	"github.com/open-feature/go-sdk/openfeature"
)

// NewClient for feature.
func NewClient(cfg *Config) *openfeature.Client {
	openfeature.SetProvider(provider(cfg))

	return openfeature.NewClient(os.ExecutableName())
}

func provider(cfg *Config) openfeature.FeatureProvider {
	if cfg.Kind == "flipt" {
		return flipt.NewProvider(flipt.WithAddress(cfg.Host))
	}

	return openfeature.NoopProvider{}
}
