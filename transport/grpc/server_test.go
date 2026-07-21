package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/config/server"
	"github.com/alexfalkowski/go-service/v2/context"
	tls "github.com/alexfalkowski/go-service/v2/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	netgrpc "github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/net/http"
	netserver "github.com/alexfalkowski/go-service/v2/net/server"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/grpc"
	"github.com/alexfalkowski/go-sync"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func init() {
	grpc.Register(test.FS)
}

func TestServer(t *testing.T) {
	mux := http.NewServeMux()
	lc := fxtest.NewLifecycle(t)
	logger, err := test.NewLogger(lc, test.NewTextLoggerConfig())
	require.NoError(t, err)

	c := test.NewInsecureTransportConfig()
	s := &test.Server{Lifecycle: lc, Logger: logger, TransportConfig: c, Mux: mux, Drain: netserver.NewDrain()}
	require.NoError(t, s.Register())

	lc.RequireStart()
	lc.RequireStop()
}

func TestDisabledServer(t *testing.T) {
	srv, err := grpc.NewServer(grpc.ServerParams{Config: nil})

	require.NoError(t, err)
	require.Nil(t, srv)
}

func TestNilServer(t *testing.T) {
	var srv *grpc.Server

	require.Nil(t, srv.ServiceRegistrar())
	require.Nil(t, srv.GetService())
}

func TestInvalidServer(t *testing.T) {
	cfg := &grpc.Config{
		Config: &server.Config{
			Timeout: 5 * time.Second,
			TLS:     test.NewTLSConfig("certs/client-cert.pem", "secrets/none"),
		},
	}
	params := grpc.ServerParams{
		Shutdowner:   test.NewShutdowner(),
		Config:       cfg,
		MethodPolicy: grpc.NewMethodPolicy(),
	}

	_, err := grpc.NewServer(params)
	require.Error(t, err)
}

func TestServerRejectsCAOnlyTLS(t *testing.T) {
	cfg := &grpc.Config{
		Config: &server.Config{
			Timeout: 5 * time.Second,
			TLS:     &tls.Config{CA: test.FilePath("certs/rootCA.pem")},
		},
	}
	params := grpc.ServerParams{
		Shutdowner:   test.NewShutdowner(),
		Config:       cfg,
		MethodPolicy: grpc.NewMethodPolicy(),
	}

	_, err := grpc.NewServer(params)
	require.ErrorIs(t, err, server.ErrMissingKeyPair)
	require.EqualError(t, err, "grpc: server: missing tls key pair")
}

func TestServerAppliesCallerServerOption(t *testing.T) {
	rejected := errors.New("rejected by tap handle")
	tapHandle := func(ctx context.Context, _ *netgrpc.TapInfo) (context.Context, error) {
		return ctx, rejected
	}

	cfg := test.NewInsecureTransportConfig().GRPC
	params := grpc.ServerParams{
		Shutdowner:   test.NewShutdowner(),
		Config:       cfg,
		MethodPolicy: grpc.NewMethodPolicy(),
		Options:      []netgrpc.ServerOption{netgrpc.InTapHandle(tapHandle)},
	}

	srv, err := grpc.NewServer(params)
	require.NoError(t, err)

	v1.RegisterGreeterServiceServer(srv.ServiceRegistrar(), test.NewService())
	srv.GetService().Start()
	t.Cleanup(func() {
		require.NoError(t, srv.GetService().Stop(context.Background()))
	})

	conn, err := netgrpc.NewClient(srv.GetService().String(), netgrpc.WithTransportCredentials(netgrpc.NewInsecureCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, conn.Close())
	})

	_, err = v1.NewGreeterServiceClient(conn).SayHello(t.Context(), &v1.SayHelloRequest{Name: "test"})

	require.Error(t, err)
	require.ErrorContains(t, err, rejected.Error())
}

func TestServerRecoversUnaryMetadataPanic(t *testing.T) {
	client := newPanicMetadataServerClient(t)

	_, err := client.SayHello(t.Context(), &v1.SayHelloRequest{Name: "panic"})
	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))

	resp, err := client.SayHello(t.Context(), &v1.SayHelloRequest{Name: "test"})
	require.NoError(t, err)
	require.Equal(t, "Hello test", resp.GetMessage())
}

func TestServerRecoversStreamMetadataPanic(t *testing.T) {
	client := newPanicMetadataServerClient(t)

	stream, err := client.SayStreamHello(t.Context())
	require.NoError(t, err)

	_, err = test.SendStreamHello(t, stream, "panic")
	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))

	stream, err = client.SayStreamHello(t.Context())
	require.NoError(t, err)

	resp, err := test.SendStreamHello(t, stream, "test")
	require.NoError(t, err)
	require.Equal(t, "Hello test", resp.GetMessage())
}

func newPanicMetadataServerClient(t *testing.T) v1.GreeterServiceClient {
	t.Helper()

	params := grpc.ServerParams{
		Shutdowner:   test.NewShutdowner(),
		Config:       test.NewInsecureTransportConfig().GRPC,
		ID:           &panicOnceIDGenerator{},
		MethodPolicy: grpc.NewMethodPolicy(),
	}

	srv, err := grpc.NewServer(params)
	require.NoError(t, err)

	v1.RegisterGreeterServiceServer(srv.ServiceRegistrar(), test.NewService())
	srv.GetService().Start()
	t.Cleanup(func() {
		require.NoError(t, srv.GetService().Stop(context.Background()))
	})

	conn, err := netgrpc.NewClient(srv.GetService().String(), netgrpc.WithTransportCredentials(netgrpc.NewInsecureCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, conn.Close())
	})

	return v1.NewGreeterServiceClient(conn)
}

var _ id.Generator = (*panicOnceIDGenerator)(nil)

type panicOnceIDGenerator struct {
	once sync.Once
}

func (g *panicOnceIDGenerator) Generate() string {
	g.once.Do(func() {
		panic("metadata panic")
	})

	return "request-id"
}
