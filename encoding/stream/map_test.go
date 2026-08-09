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

			codec := m.Get(kind)
			require.NotNil(t, codec.Encoder)
			require.NotNil(t, codec.Decoder)

			require.IsType(t, expectedEncoders[kind], codec.Encoder(io.Discard))
			require.IsType(t, expectedDecoders[kind], codec.Decoder(strings.NewReader("")))
		})
	}
}

func TestMapGetIsNilSafe(t *testing.T) {
	t.Parallel()

	var m *stream.Map

	require.Equal(t, stream.Codec{}, m.Get("json"))
}

func TestMapGetReturnsZeroCodecForUnknownKind(t *testing.T) {
	t.Parallel()

	m := stream.NewMap()

	for _, kind := range []string{"test", "bob"} {
		t.Run(kind, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, stream.Codec{}, m.Get(kind))
		})
	}
}

func TestMapKeysIncludesPartiallyRegisteredKinds(t *testing.T) {
	t.Parallel()

	m := stream.NewMap()
	m.Register("encode-only", stream.Codec{Encoder: func(w io.Writer) stream.Encoder { return json.NewEncoder(w) }})
	m.Register("decode-only", stream.Codec{Decoder: func(r io.Reader) stream.Decoder { return json.NewDecoder(r) }})

	require.ElementsMatch(t, []string{"json", "msgpack", "gob", "yaml", "encode-only", "decode-only"}, m.Keys())
}

func TestMapRegisterReplacesExistingCodec(t *testing.T) {
	t.Parallel()

	m := stream.NewMap()
	encoder := func(w io.Writer) stream.Encoder { return json.NewEncoder(w) }
	decoder := func(r io.Reader) stream.Decoder { return json.NewDecoder(r) }

	m.Register("custom", stream.Codec{Encoder: encoder, Decoder: decoder})
	codec := m.Get("custom")
	require.Equal(t, reflect.ValueOf(encoder).Pointer(), reflect.ValueOf(codec.Encoder).Pointer())
	require.Equal(t, reflect.ValueOf(decoder).Pointer(), reflect.ValueOf(codec.Decoder).Pointer())

	replacementEncoder := func(w io.Writer) stream.Encoder { return msgpack.NewEncoder(w) }
	replacementDecoder := func(r io.Reader) stream.Decoder { return yaml.NewDecoder(r) }

	m.Register("custom", stream.Codec{Encoder: replacementEncoder, Decoder: replacementDecoder})
	codec = m.Get("custom")
	require.Equal(t, reflect.ValueOf(replacementEncoder).Pointer(), reflect.ValueOf(codec.Encoder).Pointer())
	require.Equal(t, reflect.ValueOf(replacementDecoder).Pointer(), reflect.ValueOf(codec.Decoder).Pointer())

	// Register replaces the codec in full rather than merging fields: registering a partial Codec must
	// clear the side that is left unset, not keep the previous decoder around.
	anotherEncoder := func(w io.Writer) stream.Encoder { return json.NewEncoder(w) }
	m.Register("custom", stream.Codec{Encoder: anotherEncoder})
	codec = m.Get("custom")
	require.Equal(t, reflect.ValueOf(anotherEncoder).Pointer(), reflect.ValueOf(codec.Encoder).Pointer())
	require.Nil(t, codec.Decoder)
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

			codec := m.Get(kind)
			require.NotNil(t, codec.Encoder)
			require.NotNil(t, codec.Decoder)
		})
	}
}
