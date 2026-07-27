package encoding_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/encoding/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding/gob"
	"github.com/alexfalkowski/go-service/v2/encoding/hjson"
	"github.com/alexfalkowski/go-service/v2/encoding/json"
	"github.com/alexfalkowski/go-service/v2/encoding/msgpack"
	"github.com/alexfalkowski/go-service/v2/encoding/proto"
	"github.com/alexfalkowski/go-service/v2/encoding/toml"
	"github.com/alexfalkowski/go-service/v2/encoding/yaml"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestEncoder(t *testing.T) {
	for _, k := range test.Encoder.Keys() {
		t.Run(k, func(t *testing.T) {
			require.NotNil(t, test.Encoder.Get(k))
		})
	}

	for _, k := range []string{"test", "bob"} {
		t.Run(k, func(t *testing.T) {
			require.Nil(t, test.Encoder.Get(k))
		})
	}
}

func TestNewMapRegistersDefaultEncoders(t *testing.T) {
	encoders := encoding.NewMap()

	expected := map[string]encoding.Encoder{
		"json":         json.NewEncoder(),
		"hjson":        hjson.NewEncoder(),
		"yaml":         yaml.NewEncoder(),
		"yml":          yaml.NewEncoder(),
		"toml":         toml.NewEncoder(),
		"msgpack":      msgpack.NewEncoder(),
		"pb":           proto.NewBinary(),
		"pbbin":        proto.NewBinary(),
		"proto":        proto.NewBinary(),
		"protobin":     proto.NewBinary(),
		"protobuf":     proto.NewBinary(),
		"pbtxt":        proto.NewText(),
		"prototext":    proto.NewText(),
		"prototxt":     proto.NewText(),
		"protojson":    proto.NewJSON(),
		"pbjson":       proto.NewJSON(),
		"gob":          gob.NewEncoder(),
		"markdown":     bytes.NewEncoder(),
		"octet-stream": bytes.NewEncoder(),
		"plain":        bytes.NewEncoder(),
	}

	for kind, expectedEncoder := range expected {
		t.Run(kind, func(t *testing.T) {
			require.IsType(t, expectedEncoder, encoders.Get(kind))
		})
	}
}

func TestMapRegister(t *testing.T) {
	encoders := encoding.NewMap()
	custom := test.NewEncoder(test.ErrFailed)
	replacement := bytes.NewEncoder()

	encoders.Register("custom", custom)
	require.Same(t, custom, encoders.Get("custom"))

	encoders.Register("custom", replacement)
	require.Same(t, replacement, encoders.Get("custom"))
}

func TestModuleProvidesDefaultEncoders(t *testing.T) {
	var encoders *encoding.Map

	app := fx.New(
		encoding.Module,
		fx.Populate(&encoders),
		fx.NopLogger,
	)

	require.NoError(t, app.Err())
	for _, kind := range []string{"json", "hjson", "yaml", "toml", "msgpack", "proto", "prototext", "protojson", "gob", "plain"} {
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
