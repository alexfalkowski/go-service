package access_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/os"
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

	tests := []struct {
		name   string
		user   string
		system string
		action string
		access bool
	}{
		{name: "allowed read", user: "alice", system: "service", action: "read", access: true},
		{name: "denied wrong action", user: "alice", system: "service", action: "write"},
		{name: "denied wrong user", user: "bob", system: "service", action: "read"},
		{name: "denied unknown user", user: "carol", system: "service", action: "read"},
		{name: "denied unknown system", user: "alice", system: "other", action: "read"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := controller.HasAccess(tt.user, tt.system, tt.action)
			require.NoError(t, err)
			require.Equal(t, tt.access, ok)
		})
	}
}

func TestNewControllerErrors(t *testing.T) {
	tests := []struct {
		err    error
		config *access.Config
		name   string
	}{
		{
			name: "missing model env source",
			config: &access.Config{
				Model:  "env:ACCESS_MODEL",
				Policy: test.FilePath("configs/rbac.csv"),
			},
			err: os.ErrEnvSourceMissing,
		},
		{
			name: "missing policy env source",
			config: &access.Config{
				Model:  test.FilePath("configs/rbac.conf"),
				Policy: "env:ACCESS_POLICY",
			},
			err: os.ErrEnvSourceMissing,
		},
		{
			name: "invalid model content",
			config: &access.Config{
				Model:  "not a casbin model",
				Policy: test.FilePath("configs/rbac.csv"),
			},
		},
		{
			name: "empty policy content",
			config: &access.Config{
				Model: test.FilePath("configs/rbac.conf"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller, err := access.NewController(tt.config, test.FS)
			require.Error(t, err)
			require.Nil(t, controller)

			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
			}
		})
	}
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

	ok, err := controller.HasAccess("alice", "service", "read")
	require.NoError(t, err)
	require.True(t, ok)
}

func TestHasAccessWithLiteralSources(t *testing.T) {
	model, err := test.FS.ReadFile(test.Path("configs/rbac.conf"))
	require.NoError(t, err)

	policy, err := test.FS.ReadFile(test.Path("configs/rbac.csv"))
	require.NoError(t, err)

	controller, err := access.NewController(&access.Config{
		Model:  string(model),
		Policy: string(policy),
	}, test.FS)
	require.NoError(t, err)

	ok, err := controller.HasAccess("alice", "service", "read")
	require.NoError(t, err)
	require.True(t, ok)
}
