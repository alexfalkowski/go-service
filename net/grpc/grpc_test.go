package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/config/options"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestStatusText(t *testing.T) {
	require.Equal(t, codes.Unauthenticated.String(), grpc.StatusText(codes.Unauthenticated))
}

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

func TestNewServerRejectsNegativeTimeoutOption(t *testing.T) {
	keys := []string{
		"keepalive_enforcement_policy_ping_min_time",
		"keepalive_max_connection_idle",
		"keepalive_max_connection_age",
		"keepalive_max_connection_age_grace",
		"keepalive_ping_time",
		"connection_timeout",
	}

	for _, key := range keys {
		t.Run(key, func(t *testing.T) {
			require.Panics(t, func() {
				grpc.NewServer(options.Map{key: "-1s"}, time.Second)
			})
		})
	}
}

func TestNewServerWithOverflowingAdvancedOptions(t *testing.T) {
	require.Panics(t, func() {
		grpc.NewServer(options.Map{"initial_window_size": "3GB"}, time.Second)
	})

	require.Panics(t, func() {
		grpc.NewServer(options.Map{"max_header_list_size": "5GB"}, time.Second)
	})
}
