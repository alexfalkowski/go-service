package access_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/token/access"
	"github.com/stretchr/testify/require"
)

func TestNewController(t *testing.T) {
	_, err := access.NewController(&access.Config{
		Model:  test.FilePath("configs/rbac.conf"),
		Policy: test.FilePath("configs/bob"),
	}, test.FS)
	require.Error(t, err)

	controller, err := access.NewController(nil, test.FS)
	require.NoError(t, err)
	require.Nil(t, controller)
}

func TestHasAccess(t *testing.T) {
	config := test.NewAccessConfig()

	controller, err := access.NewController(config, test.FS)
	require.NoError(t, err)

	ok, err := controller.HasAccess("alice", "service:read")
	require.NoError(t, err)
	require.True(t, ok)
}

func TestHasAccessWithEnvSources(t *testing.T) {
	model, err := test.FS.ReadFile(test.Path("configs/rbac.conf"))
	require.NoError(t, err)

	policy, err := test.FS.ReadFile(test.Path("configs/rbac.csv"))
	require.NoError(t, err)

	t.Setenv("ACCESS_MODEL", string(model))
	t.Setenv("ACCESS_POLICY", string(policy))

	controller, err := access.NewController(&access.Config{
		Model:  "env:ACCESS_MODEL",
		Policy: "env:ACCESS_POLICY",
	}, test.FS)
	require.NoError(t, err)

	ok, err := controller.HasAccess("alice", "service:read")
	require.NoError(t, err)
	require.True(t, ok)
}
