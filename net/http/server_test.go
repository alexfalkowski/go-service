package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/config/options"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/stretchr/testify/require"
)

func TestNewServerRejectsNegativeTimeoutOption(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})

	for _, key := range []string{"read_timeout", "write_timeout", "idle_timeout", "read_header_timeout"} {
		t.Run(key, func(t *testing.T) {
			t.Parallel()

			require.Panics(t, func() {
				http.NewServer(options.Map{key: "-1s"}, handler)
			})
		})
	}
}

func TestNewServerSetsProtocols(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})

	server := http.NewServer(options.Map{}, handler)

	require.Nil(t, server.HTTP2)
	require.NotNil(t, server.Protocols)
	require.True(t, server.Protocols.HTTP1())
	require.True(t, server.Protocols.HTTP2())
	require.True(t, server.Protocols.UnencryptedHTTP2())
}

func TestNewServerSetsHTTP2Options(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                       string
		options                    options.Map
		maxConcurrentStreams       int
		maxReceiveBufferConnection int
		maxReceiveBufferStream     int
	}{
		{
			name:                 "max concurrent streams",
			options:              options.Map{"http2_max_concurrent_streams": "7"},
			maxConcurrentStreams: 7,
		},
		{
			name:                       "max receive buffer per connection",
			options:                    options.Map{"http2_max_receive_buffer_per_connection": "3MB"},
			maxReceiveBufferConnection: int(3 * bytes.MB.Bytes()),
		},
		{
			name:                   "max receive buffer per stream",
			options:                options.Map{"http2_max_receive_buffer_per_stream": "4MB"},
			maxReceiveBufferStream: int(4 * bytes.MB.Bytes()),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			server := http.NewServer(test.options, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))

			require.NotNil(t, server.HTTP2)
			require.Equal(t, test.maxConcurrentStreams, server.HTTP2.MaxConcurrentStreams)
			require.Equal(t, test.maxReceiveBufferConnection, server.HTTP2.MaxReceiveBufferPerConnection)
			require.Equal(t, test.maxReceiveBufferStream, server.HTTP2.MaxReceiveBufferPerStream)
		})
	}
}

func TestNewServerRejectsOverflowingHTTP2Options(t *testing.T) {
	t.Parallel()

	for key, value := range map[string]string{
		"http2_max_concurrent_streams":            "4294967296",
		"http2_max_receive_buffer_per_connection": "3GB",
		"http2_max_receive_buffer_per_stream":     "3GB",
	} {
		t.Run(key, func(t *testing.T) {
			t.Parallel()

			require.Panics(t, func() {
				http.NewServer(options.Map{key: value}, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
			})
		})
	}
}
