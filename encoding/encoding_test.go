package encoding_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/encoding/json"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestMapReturnsRegisteredEncoders(t *testing.T) {
	for _, k := range test.Encoder.Keys() {
		t.Run(k, func(t *testing.T) {
			require.NotNil(t, test.Encoder.Get(k))
		})
	}

	for _, k := range []string{"test", "bob"} {
		t.Run("returns nil for "+k, func(t *testing.T) {
			require.Nil(t, test.Encoder.Get(k))
		})
	}
}

func TestModuleProvidesDefaultEncoders(t *testing.T) {
	var encoders *encoding.Map

	app := fx.New(
		encoding.Module,
		fx.Populate(&encoders),
		fx.NopLogger,
	)

	require.NoError(t, app.Err())
	for _, kind := range []string{"json", "hjson", "yaml", "toml", "msgpack", "protobuf", "prototext", "protojson", "gob", "bytes"} {
		t.Run(kind, func(t *testing.T) {
			require.NotNil(t, encoders.Get(kind))
		})
	}
}

func TestModuleDoesNotProvideJSONEncoder(t *testing.T) {
	var encoder *json.Encoder

	app := fx.New(
		encoding.Module,
		fx.Populate(&encoder),
		fx.NopLogger,
	)

	require.Error(t, app.Err())
}
