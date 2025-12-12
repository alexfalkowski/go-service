package env_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	require.Equal(t, "(devel)", env.NewVersion().String())
	require.Equal(t, "1.0.0", env.Version("v1.0.0").String())
	require.Equal(t, "what", env.Version("what").String())
	require.Empty(t, env.Version(strings.Empty).String())

	t.Setenv("SERVICE_VERSION", test.Version.String())
	require.Equal(t, "1.0.0", env.NewVersion().String())
}
