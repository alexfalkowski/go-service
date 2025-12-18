package env_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestName(t *testing.T) {
	require.Equal(t, "env.test", env.NewName(test.FS).String())

	t.Setenv("SERVICE_NAME", test.Name.String())
	require.Equal(t, "test", env.NewName(test.FS).String())
}

func TestUserAgent(t *testing.T) {
	require.Equal(t, "test/1.0.0", env.NewUserAgent(env.Name("test"), env.Version("v1.0.0")).String())
	require.Equal(t, "test/test", env.NewUserAgent(env.Name("test"), env.Version("test")).String())
}

func TestID(t *testing.T) {
	generator := uuid.NewGenerator()
	require.NotEmpty(t, env.NewID(generator).String())

	t.Setenv("SERVICE_ID", "new_id")
	require.Equal(t, "new_id", env.NewID(generator).String())
}
