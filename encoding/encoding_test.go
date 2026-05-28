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
	jsonEncoder := json.NewEncoder()
	hjsonEncoder := hjson.NewEncoder()
	yamlEncoder := yaml.NewEncoder()
	tomlEncoder := toml.NewEncoder()
	msgpackEncoder := msgpack.NewEncoder()
	protoBinary := proto.NewBinary()
	protoText := proto.NewText()
	protoJSON := proto.NewJSON()
	gobEncoder := gob.NewEncoder()
	bytesEncoder := bytes.NewEncoder()

	encoders := encoding.NewMap(encoding.MapParams{
		JSON:        jsonEncoder,
		HumanJSON:   hjsonEncoder,
		YAML:        yamlEncoder,
		TOML:        tomlEncoder,
		MessagePack: msgpackEncoder,
		ProtoBinary: protoBinary,
		ProtoText:   protoText,
		ProtoJSON:   protoJSON,
		GOB:         gobEncoder,
		Bytes:       bytesEncoder,
	})

	expected := map[string]encoding.Encoder{
		"json":         jsonEncoder,
		"hjson":        hjsonEncoder,
		"yaml":         yamlEncoder,
		"yml":          yamlEncoder,
		"toml":         tomlEncoder,
		"msgpack":      msgpackEncoder,
		"pb":           protoBinary,
		"pbbin":        protoBinary,
		"proto":        protoBinary,
		"protobin":     protoBinary,
		"protobuf":     protoBinary,
		"pbtxt":        protoText,
		"prototext":    protoText,
		"prototxt":     protoText,
		"protojson":    protoJSON,
		"pbjson":       protoJSON,
		"gob":          gobEncoder,
		"markdown":     bytesEncoder,
		"octet-stream": bytesEncoder,
		"plain":        bytesEncoder,
	}

	for kind, expectedEncoder := range expected {
		t.Run(kind, func(t *testing.T) {
			require.Same(t, expectedEncoder, encoders.Get(kind))
		})
	}
}

func TestMapRegister(t *testing.T) {
	encoders := encoding.NewMap(encoding.MapParams{})
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
