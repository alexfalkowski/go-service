package options_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/config/options"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestDuration(t *testing.T) {
	require.Equal(t, 10*time.Minute, test.ConfigOptions.Duration("read_timeout", 5*time.Second))
	require.Equal(t, 5*time.Second, test.ConfigOptions.Duration("timeout", 5*time.Second))
	require.Equal(t, 5*time.Second, test.ConfigOptions.Duration("bob", 5*time.Second))
}

func TestUint32(t *testing.T) {
	opts := options.Map{"count": "12"}

	require.EqualValues(t, 12, opts.Uint32("count", 5))
	require.EqualValues(t, 5, opts.Uint32("missing", 5))
	require.Panics(t, func() { options.Map{"negative": "-1"}.Uint32("negative", 5) })
	require.Panics(t, func() { options.Map{"invalid": "abc"}.Uint32("invalid", 5) })
	require.Panics(t, func() { options.Map{"overflow": "4294967296"}.Uint32("overflow", 5) })
}

func TestSize(t *testing.T) {
	opts := options.Map{"limit": "16MB"}

	require.Equal(t, 16*bytes.MB, opts.Size("limit", bytes.MB))
	require.Equal(t, bytes.MB, opts.Size("missing", bytes.MB))
	require.Panics(t, func() { options.Map{"invalid": "abc"}.Size("invalid", bytes.MB) })
}
