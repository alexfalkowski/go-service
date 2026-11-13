package encoding_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v3/encoding"
	"github.com/alexfalkowski/go-service/v3/encoding/json"
	"github.com/alexfalkowski/go-service/v3/internal/test"
	"github.com/alexfalkowski/go-signal/v2"
	"github.com/stretchr/testify/require"
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

	_, err := signal.NewApp(
		encoding.Module,
		signal.Register(func(value *encoding.Map) { encoders = value }),
	)

	require.NoError(t, err)
	for _, kind := range []string{"json", "hjson", "yaml", "toml", "msgpack", "protobuf", "prototext", "protojson", "gob", "bytes"} {
		t.Run(kind, func(t *testing.T) {
			require.NotNil(t, encoders.Get(kind))
		})
	}
}

func TestModuleDoesNotProvideJSONEncoder(t *testing.T) {
	_, err := signal.NewApp(
		encoding.Module,
		signal.Register(func(*json.Encoder) {}),
	)

	require.Error(t, err)
}
