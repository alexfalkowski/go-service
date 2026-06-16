package transport_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/transport"
	"github.com/stretchr/testify/require"
)

func TestNewAccessControllerWithoutAccessConfig(t *testing.T) {
	controller, err := transport.NewAccessController(nil, test.FS)
	require.NoError(t, err)
	require.Nil(t, controller)
}

func TestNewAccessControllerWithAccessConfig(t *testing.T) {
	controller, err := transport.NewAccessController(test.NewAccessConfig(), test.FS)
	require.NoError(t, err)
	require.NotNil(t, controller)
}
