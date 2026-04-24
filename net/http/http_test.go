package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/config/options"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestNewServerWithMaxHeaderBytes(t *testing.T) {
	server := http.NewServer(options.Map{"max_header_bytes": "16KB"}, time.Second, http.NewServeMux())

	require.Equal(t, int(16*bytes.KB), server.MaxHeaderBytes)
}

func TestNewServerWithDefaultMaxHeaderBytes(t *testing.T) {
	server := http.NewServer(options.Map{}, time.Second, http.NewServeMux())

	require.Equal(t, http.DefaultMaxHeaderBytes, server.MaxHeaderBytes)
}
