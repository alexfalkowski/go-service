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

	t.Setenv("SERVICE_VERSION", test.Version.String())
	require.Equal(t, "1.0.0", env.NewVersion().String())
}

func TestVersionString(t *testing.T) {
	for _, tt := range []struct {
		name     string
		version  env.Version
		expected string
	}{
		{name: "semver prefix", version: "v1.0.0", expected: "1.0.0"},
		{name: "plain value", version: "what", expected: "what"},
		{name: "empty", version: env.Version(strings.Empty), expected: strings.Empty},
	} {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.version.String())
		})
	}
}
