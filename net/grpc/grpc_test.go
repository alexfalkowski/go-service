package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/config/options"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestNewServerWithAdvancedOptions(t *testing.T) {
	opts := options.Map{
		"max_concurrent_streams":   "7",
		"connection_timeout":       "250ms",
		"max_header_list_size":     "9MB",
		"initial_window_size":      "3MB",
		"initial_conn_window_size": "4MB",
		"max_send_msg_size":        "5MB",
	}

	require.NotPanics(t, func() {
		server := grpc.NewServer(opts, time.Second)
		require.NotNil(t, server)
	})
}

func TestNewServerWithOverflowingAdvancedOptions(t *testing.T) {
	require.Panics(t, func() {
		grpc.NewServer(options.Map{"initial_window_size": "3GB"}, time.Second)
	})

	require.Panics(t, func() {
		grpc.NewServer(options.Map{"max_header_list_size": "5GB"}, time.Second)
	})
}
