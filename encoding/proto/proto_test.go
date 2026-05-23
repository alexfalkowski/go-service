package proto_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/encoding/errors"
	"github.com/alexfalkowski/go-service/v2/encoding/proto"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/grpc/health"
	"github.com/stretchr/testify/require"
)

func TestValidBinaryEncoder(t *testing.T) {
	encoder := proto.NewBinary()

	bytes := test.Pool.Get()
	defer test.Pool.Put(bytes)

	require.NoError(t, encoder.Encode(bytes, &health.Response{Status: health.Serving}))

	var decode health.Response
	require.NoError(t, encoder.Decode(bytes, &decode))
	require.Equal(t, health.Serving, decode.GetStatus())
}

func TestInvalidBinaryEncode(t *testing.T) {
	encoder := proto.NewBinary()

	bytes := test.Pool.Get()
	defer test.Pool.Put(bytes)

	var msg string
	require.Error(t, encoder.Encode(bytes, &msg))
}

func TestInvalidBinaryDecode(t *testing.T) {
	encoder := proto.NewBinary()

	bytes := test.Pool.Get()
	defer test.Pool.Put(bytes)

	require.NoError(t, encoder.Encode(bytes, &health.Response{Status: health.Serving}))

	var msg string
	require.Error(t, encoder.Decode(bytes, &msg))
}

func TestInvalidBinaryDecodeDoesNotRead(t *testing.T) {
	encoder := proto.NewBinary()
	reader := &test.TrackingReader{Err: io.EOF}

	var msg string
	err := encoder.Decode(reader, &msg)
	require.ErrorIs(t, err, errors.ErrInvalidType)
	require.Zero(t, reader.Reads)
}

func TestValidTextEncoder(t *testing.T) {
	encoder := proto.NewText()

	bytes := test.Pool.Get()
	defer test.Pool.Put(bytes)

	require.NoError(t, encoder.Encode(bytes, &health.Response{Status: health.Serving}))

	var decode health.Response
	require.NoError(t, encoder.Decode(bytes, &decode))
	require.Equal(t, health.Serving, decode.GetStatus())
}

func TestInvalidTextEncode(t *testing.T) {
	encoder := proto.NewText()

	bytes := test.Pool.Get()
	defer test.Pool.Put(bytes)

	var msg string
	require.Error(t, encoder.Encode(bytes, &msg))
}

func TestInvalidTextDecode(t *testing.T) {
	encoder := proto.NewText()

	bytes := test.Pool.Get()
	defer test.Pool.Put(bytes)

	require.NoError(t, encoder.Encode(bytes, &health.Response{Status: health.Serving}))

	var msg string
	require.Error(t, encoder.Decode(bytes, &msg))
}

func TestInvalidTextDecodeDoesNotRead(t *testing.T) {
	encoder := proto.NewText()
	reader := &test.TrackingReader{Err: io.EOF}

	var msg string
	err := encoder.Decode(reader, &msg)
	require.ErrorIs(t, err, errors.ErrInvalidType)
	require.Zero(t, reader.Reads)
}

func TestValidJSONEncoder(t *testing.T) {
	encoder := proto.NewJSON()

	bytes := test.Pool.Get()
	defer test.Pool.Put(bytes)

	require.NoError(t, encoder.Encode(bytes, &health.Response{Status: health.Serving}))

	var decode health.Response
	require.NoError(t, encoder.Decode(bytes, &decode))
	require.Equal(t, health.Serving, decode.GetStatus())
}

func TestInvalidJSONEncode(t *testing.T) {
	encoder := proto.NewJSON()

	bytes := test.Pool.Get()
	defer test.Pool.Put(bytes)

	var msg string
	require.Error(t, encoder.Encode(bytes, &msg))
}

func TestInvalidJSONDecode(t *testing.T) {
	encoder := proto.NewJSON()

	bytes := test.Pool.Get()
	defer test.Pool.Put(bytes)

	require.NoError(t, encoder.Encode(bytes, &health.Response{Status: health.Serving}))

	var msg string
	require.Error(t, encoder.Decode(bytes, &msg))
}

func TestInvalidJSONDecodeDoesNotRead(t *testing.T) {
	encoder := proto.NewJSON()
	reader := &test.TrackingReader{Err: io.EOF}

	var msg string
	err := encoder.Decode(reader, &msg)
	require.ErrorIs(t, err, errors.ErrInvalidType)
	require.Zero(t, reader.Reads)
}

func TestEncodeReturnsWriteError(t *testing.T) {
	encoders := []struct {
		encoder encoding.Encoder
		name    string
	}{
		{encoder: proto.NewBinary(), name: "binary"},
		{encoder: proto.NewJSON(), name: "json"},
		{encoder: proto.NewText(), name: "text"},
	}

	for _, tt := range encoders {
		t.Run(tt.name, func(t *testing.T) {
			require.Error(t, tt.encoder.Encode(test.ErrWriter{}, &health.Response{Status: health.Serving}))
		})
	}
}

func TestErrBinaryDecode(t *testing.T) {
	encoder := proto.NewBinary()
	var msg health.Response
	require.Error(t, encoder.Decode(&test.ErrReaderCloser{}, &msg))
}

func TestErrTextDecode(t *testing.T) {
	encoder := proto.NewText()
	var msg health.Response
	require.Error(t, encoder.Decode(&test.ErrReaderCloser{}, &msg))
}

func TestErrJSONDecode(t *testing.T) {
	encoder := proto.NewJSON()
	var msg health.Response
	require.Error(t, encoder.Decode(&test.ErrReaderCloser{}, &msg))
}

func TestInvalidTypedNilEncode(t *testing.T) {
	encoders := []struct {
		encoder encoding.Encoder
		name    string
	}{
		{encoder: proto.NewBinary(), name: "binary"},
		{encoder: proto.NewJSON(), name: "json"},
		{encoder: proto.NewText(), name: "text"},
	}

	for _, tt := range encoders {
		t.Run(tt.name, func(t *testing.T) {
			var msg *health.Response

			err := tt.encoder.Encode(io.Discard, msg)
			require.ErrorIs(t, err, errors.ErrInvalidType)
		})
	}
}

func TestInvalidTypedNilDecodeDoesNotRead(t *testing.T) {
	encoders := []struct {
		encoder encoding.Encoder
		name    string
	}{
		{encoder: proto.NewBinary(), name: "binary"},
		{encoder: proto.NewJSON(), name: "json"},
		{encoder: proto.NewText(), name: "text"},
	}

	for _, tt := range encoders {
		t.Run(tt.name, func(t *testing.T) {
			reader := &test.TrackingReader{Err: io.EOF}
			var msg *health.Response

			err := tt.encoder.Decode(reader, msg)
			require.ErrorIs(t, err, errors.ErrInvalidType)
			require.Zero(t, reader.Reads)
		})
	}
}
