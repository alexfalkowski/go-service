package proto_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/encoding/errors"
	"github.com/alexfalkowski/go-service/v2/encoding/proto"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/grpc/health"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestEncodeDecode(t *testing.T) {
	for _, tt := range protoEncoders() {
		t.Run(tt.name, func(t *testing.T) {
			buffer := test.Pool.Get()
			defer test.Pool.Put(buffer)

			require.NoError(t, tt.encoder.Encode(buffer, &health.Response{Status: health.Serving}))

			var decode health.Response
			require.NoError(t, tt.encoder.Decode(buffer, &decode))
			require.Equal(t, health.Serving, decode.GetStatus())
		})
	}
}

func TestEncodeRejectsInvalidType(t *testing.T) {
	for _, tt := range protoEncoders() {
		t.Run(tt.name, func(t *testing.T) {
			buffer := test.Pool.Get()
			defer test.Pool.Put(buffer)

			var msg string
			require.ErrorIs(t, tt.encoder.Encode(buffer, &msg), errors.ErrInvalidType)
		})
	}
}

func TestDecodeRejectsInvalidType(t *testing.T) {
	for _, tt := range protoEncoders() {
		t.Run(tt.name, func(t *testing.T) {
			buffer := test.Pool.Get()
			defer test.Pool.Put(buffer)

			require.NoError(t, tt.encoder.Encode(buffer, &health.Response{Status: health.Serving}))

			var msg string
			require.ErrorIs(t, tt.encoder.Decode(buffer, &msg), errors.ErrInvalidType)
		})
	}
}

func TestDecodeRejectsInvalidTypeWithoutReading(t *testing.T) {
	for _, tt := range protoEncoders() {
		t.Run(tt.name, func(t *testing.T) {
			reader := &test.TrackingReader{Err: io.EOF}

			var msg string
			err := tt.encoder.Decode(reader, &msg)
			require.ErrorIs(t, err, errors.ErrInvalidType)
			require.Zero(t, reader.Reads)
		})
	}
}

func TestEncodersUseExpectedWireFormats(t *testing.T) {
	t.Run("binary", func(t *testing.T) {
		encoder := proto.NewBinary()
		buffer := test.Pool.Get()
		defer test.Pool.Put(buffer)

		require.NoError(t, encoder.Encode(buffer, &health.Response{Status: health.Serving}))
		require.Equal(t, []byte{0x08, 0x01}, test.Pool.Copy(buffer))

		var decode health.Response
		require.NoError(t, encoder.Decode(bytes.NewBuffer([]byte{0x08, 0x01}), &decode))
		require.Equal(t, health.Serving, decode.GetStatus())
	})

	t.Run("text", func(t *testing.T) {
		encoder := proto.NewText()
		buffer := test.Pool.Get()
		defer test.Pool.Put(buffer)

		require.NoError(t, encoder.Encode(buffer, &health.Response{Status: health.Serving}))
		require.Equal(t, "status:SERVING", strings.TrimSpace(buffer.String()))

		var decode health.Response
		require.NoError(t, encoder.Decode(bytes.NewBufferString("status:SERVING"), &decode))
		require.Equal(t, health.Serving, decode.GetStatus())
	})

	t.Run("json", func(t *testing.T) {
		encoder := proto.NewJSON()
		buffer := test.Pool.Get()
		defer test.Pool.Put(buffer)

		require.NoError(t, encoder.Encode(buffer, &health.Response{Status: health.Serving}))
		require.JSONEq(t, `{"status":"SERVING"}`, buffer.String())

		var decode health.Response
		require.NoError(t, encoder.Decode(bytes.NewBufferString(`{"status":"SERVING"}`), &decode))
		require.Equal(t, health.Serving, decode.GetStatus())
	})
}

func TestEncodeReturnsWriteError(t *testing.T) {
	for _, tt := range protoEncoders() {
		t.Run(tt.name, func(t *testing.T) {
			require.ErrorIs(t, tt.encoder.Encode(test.ErrWriter{}, &health.Response{Status: health.Serving}), test.ErrFailed)
		})
	}
}

func TestDecodeReturnsReadError(t *testing.T) {
	for _, tt := range protoEncoders() {
		t.Run(tt.name, func(t *testing.T) {
			var msg health.Response
			require.ErrorIs(t, tt.encoder.Decode(&test.ErrReaderCloser{}, &msg), test.ErrFailed)
		})
	}
}

func TestInvalidTypedNilEncode(t *testing.T) {
	for _, tt := range protoEncoders() {
		t.Run(tt.name, func(t *testing.T) {
			var msg *health.Response

			err := tt.encoder.Encode(io.Discard, msg)
			require.ErrorIs(t, err, errors.ErrInvalidType)
		})
	}
}

func TestInvalidTypedNilDecodeDoesNotRead(t *testing.T) {
	for _, tt := range protoEncoders() {
		t.Run(tt.name, func(t *testing.T) {
			reader := &test.TrackingReader{Err: io.EOF}
			var msg *health.Response

			err := tt.encoder.Decode(reader, msg)
			require.ErrorIs(t, err, errors.ErrInvalidType)
			require.Zero(t, reader.Reads)
		})
	}
}

func protoEncoders() []struct {
	encoder encoding.Encoder
	name    string
} {
	return []struct {
		encoder encoding.Encoder
		name    string
	}{
		{encoder: proto.NewBinary(), name: "binary"},
		{encoder: proto.NewJSON(), name: "json"},
		{encoder: proto.NewText(), name: "text"},
	}
}
