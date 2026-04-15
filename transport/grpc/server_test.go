package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/config/server"
	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
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
	meter, err := test.NewPrometheusMeter(lc)
	require.NoError(t, err)

	c := test.NewInsecureTransportConfig()
	c.GRPC.TLS = &tls.Config{}

	s := &test.Server{Lifecycle: lc, Logger: logger, TransportConfig: c, Meter: meter, Mux: mux}
	require.NoError(t, s.Register())

	lc.RequireStart()
	lc.RequireStop()
}

func TestInvalidServer(t *testing.T) {
	cfg := &grpc.Config{
		Config: &server.Config{
			Timeout: 5 * time.Second,
			TLS:     test.NewTLSConfig("certs/client-cert.pem", "secrets/none"),
		},
	}
	params := grpc.ServerParams{
		Shutdowner: test.NewShutdowner(),
		Config:     cfg,
	}

	_, err := grpc.NewServer(params)
	require.Error(t, err)
}
