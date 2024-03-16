package marshaller

import (
	"errors"

	"go.uber.org/fx"
)

// ErrInvalidKind for marshaller.
var ErrInvalidKind = errors.New("invalid kind")

// Factory of marshaller.
type Factory struct {
	json *JSON
	toml *TOML
	yaml *YAML
}

// FactoryParams for marshaller.
type FactoryParams struct {
	fx.In

	JSON *JSON
	TOML *TOML
	YAML *YAML
}

// NewFactory for marshaller.
func NewFactory(params FactoryParams) *Factory {
	return &Factory{json: params.JSON, toml: params.TOML, yaml: params.YAML}
}

// Create from kind.
func (f *Factory) Create(kind string) (Marshaller, error) {
	switch kind {
	case "json":
		return f.json, nil
	case "yaml", "yml":
		return f.yaml, nil
	case "toml":
		return f.toml, nil
	}

	return nil, ErrInvalidKind
}
