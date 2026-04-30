package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/transport/http"
	"github.com/stretchr/testify/require"
)

func TestNewControllerWithoutTokenConfig(t *testing.T) {
	controller, err := http.NewController(test.NewHTTPTransportConfig(), test.FS)
	require.NoError(t, err)
	require.Nil(t, controller)
}

func TestNewTokenWithoutTokenConfig(t *testing.T) {
	tkn := http.NewToken(test.Name, test.NewHTTPTransportConfig(), nil, nil, nil, nil)
	require.Nil(t, tkn)
}
