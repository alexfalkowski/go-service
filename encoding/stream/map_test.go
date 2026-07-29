package stream_test

import (
	"reflect"
	"testing"

	"github.com/alexfalkowski/go-service/v2/encoding/stream"
	"github.com/alexfalkowski/go-service/v2/encoding/stream/gob"
	"github.com/alexfalkowski/go-service/v2/encoding/stream/json"
	"github.com/alexfalkowski/go-service/v2/encoding/stream/msgpack"
	"github.com/alexfalkowski/go-service/v2/encoding/stream/yaml"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestNewMapRegistersDefaultEncodersAndDecoders(t *testing.T) {
	t.Parallel()

	expectedEncoders := map[string]stream.Encoder{
		"json":    &json.Encoder{},
		"msgpack": &msgpack.Encoder{},
		"gob":     &gob.Encoder{},
		"yaml":    &yaml.Encoder{},
	}
	expectedDecoders := map[string]stream.Decoder{
		"json":    &json.Decoder{},
		"msgpack": &msgpack.Decoder{},
		"gob":     &gob.Decoder{},
		"yaml":    &yaml.Decoder{},
	}

	m := stream.NewMap()

	for _, kind := range []string{"json", "msgpack", "gob", "yaml"} {
		t.Run(kind, func(t *testing.T) {
			t.Parallel()

			newEncoder := m.GetEncoder(kind)
			newDecoder := m.GetDecoder(kind)
			require.NotNil(t, newEncoder)
			require.NotNil(t, newDecoder)

			require.IsType(t, expectedEncoders[kind], newEncoder(io.Discard))
			require.IsType(t, expectedDecoders[kind], newDecoder(strings.NewReader("")))
		})
	}
}

func TestMapGetIsNilSafe(t *testing.T) {
	t.Parallel()

	var m *stream.Map

	require.Nil(t, m.GetEncoder("json"))
	require.Nil(t, m.GetDecoder("json"))
}

func TestMapGetReturnsNilForUnknownKind(t *testing.T) {
	t.Parallel()

	m := stream.NewMap()

	for _, kind := range []string{"test", "bob"} {
		t.Run(kind, func(t *testing.T) {
			t.Parallel()

			require.Nil(t, m.GetEncoder(kind))
			require.Nil(t, m.GetDecoder(kind))
		})
	}
}

func TestMapKeysIsUnionOfEncodersAndDecoders(t *testing.T) {
	t.Parallel()

	m := stream.NewMap()
	m.RegisterEncoder("encode-only", func(w io.Writer) stream.Encoder { return json.NewEncoder(w) })
	m.RegisterDecoder("decode-only", func(r io.Reader) stream.Decoder { return json.NewDecoder(r) })

	require.ElementsMatch(t, []string{"json", "msgpack", "gob", "yaml", "encode-only", "decode-only"}, m.Keys())
}

func TestMapRegisterEncoder(t *testing.T) {
	t.Parallel()

	m := stream.NewMap()
	custom := func(w io.Writer) stream.Encoder { return json.NewEncoder(w) }

	m.RegisterEncoder("custom", custom)
	require.Equal(t, reflect.ValueOf(custom).Pointer(), reflect.ValueOf(m.GetEncoder("custom")).Pointer())

	replacement := func(w io.Writer) stream.Encoder { return msgpack.NewEncoder(w) }
	m.RegisterEncoder("custom", replacement)
	require.Equal(t, reflect.ValueOf(replacement).Pointer(), reflect.ValueOf(m.GetEncoder("custom")).Pointer())
}

func TestMapRegisterDecoder(t *testing.T) {
	t.Parallel()

	m := stream.NewMap()
	custom := func(r io.Reader) stream.Decoder { return json.NewDecoder(r) }

	m.RegisterDecoder("custom", custom)
	require.Equal(t, reflect.ValueOf(custom).Pointer(), reflect.ValueOf(m.GetDecoder("custom")).Pointer())

	replacement := func(r io.Reader) stream.Decoder { return yaml.NewDecoder(r) }
	m.RegisterDecoder("custom", replacement)
	require.Equal(t, reflect.ValueOf(replacement).Pointer(), reflect.ValueOf(m.GetDecoder("custom")).Pointer())
}

func TestModuleProvidesDefaultMap(t *testing.T) {
	t.Parallel()

	var m *stream.Map

	app := fx.New(
		stream.Module,
		fx.Populate(&m),
		fx.NopLogger,
	)

	require.NoError(t, app.Err())
	for _, kind := range []string{"json", "msgpack", "gob", "yaml"} {
		t.Run(kind, func(t *testing.T) {
			t.Parallel()

			require.NotNil(t, m.GetEncoder(kind))
			require.NotNil(t, m.GetDecoder(kind))
		})
	}
}
