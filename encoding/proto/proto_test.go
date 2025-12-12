package proto_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/encoding/proto"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func TestValidBinaryEncoder(t *testing.T) {
	encoder := proto.NewBinary()

	bytes := test.Pool.Get()
	defer test.Pool.Put(bytes)

	require.NoError(t, encoder.Encode(bytes, &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}))

	var decode grpc_health_v1.HealthCheckResponse
	require.NoError(t, encoder.Decode(bytes, &decode))
	require.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, decode.GetStatus())
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

	require.NoError(t, encoder.Encode(bytes, &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}))

	var msg string
	require.Error(t, encoder.Decode(bytes, &msg))
}

func TestValidTextEncoder(t *testing.T) {
	encoder := proto.NewText()

	bytes := test.Pool.Get()
	defer test.Pool.Put(bytes)

	require.NoError(t, encoder.Encode(bytes, &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}))

	var decode grpc_health_v1.HealthCheckResponse
	require.NoError(t, encoder.Decode(bytes, &decode))
	require.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, decode.GetStatus())
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

	require.NoError(t, encoder.Encode(bytes, &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}))

	var msg string
	require.Error(t, encoder.Decode(bytes, &msg))
}

func TestValidJSONEncoder(t *testing.T) {
	encoder := proto.NewJSON()

	bytes := test.Pool.Get()
	defer test.Pool.Put(bytes)

	require.NoError(t, encoder.Encode(bytes, &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}))

	var decode grpc_health_v1.HealthCheckResponse
	require.NoError(t, encoder.Decode(bytes, &decode))
	require.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, decode.GetStatus())
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

	require.NoError(t, encoder.Encode(bytes, &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}))

	var msg string
	require.Error(t, encoder.Decode(bytes, &msg))
}

func TestErrBinaryDecode(t *testing.T) {
	encoder := proto.NewBinary()
	var msg grpc_health_v1.HealthCheckResponse
	require.Error(t, encoder.Decode(&test.ErrReaderCloser{}, &msg))
}

func TestErrTextDecode(t *testing.T) {
	encoder := proto.NewText()
	var msg grpc_health_v1.HealthCheckResponse
	require.Error(t, encoder.Decode(&test.ErrReaderCloser{}, &msg))
}

func TestErrJSONDecode(t *testing.T) {
	encoder := proto.NewJSON()
	var msg grpc_health_v1.HealthCheckResponse
	require.Error(t, encoder.Decode(&test.ErrReaderCloser{}, &msg))
}
