package access_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/token/access"
	"github.com/stretchr/testify/require"
)

func TestNewController(t *testing.T) {
	_, err := access.NewController(&access.Config{Policy: test.Path("configs/bob")})
	require.Error(t, err)

	controller, err := access.NewController(nil)
	require.NoError(t, err)
	require.Nil(t, controller)
}

func TestHasAccess(t *testing.T) {
	config := test.NewAccessConfig()

	controller, err := access.NewController(config)
	require.NoError(t, err)

	ok, err := controller.HasAccess("alice", "service:read")
	require.NoError(t, err)
	require.True(t, ok)
}
