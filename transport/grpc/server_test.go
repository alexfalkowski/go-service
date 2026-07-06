package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/config/server"
	tls "github.com/alexfalkowski/go-service/v2/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	netserver "github.com/alexfalkowski/go-service/v2/net/server"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/grpc"
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
}
