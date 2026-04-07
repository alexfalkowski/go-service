package test_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/stretchr/testify/require"
)

func TestResponseWithBodyError(t *testing.T) {
	world := test.NewWorld(t)

	_, _, err := world.ResponseWithBody(t.Context(), "://bad", http.MethodGet, http.Header{}, http.NoBody)
	require.Error(t, err)
}

func TestResponseWithNoBodyError(t *testing.T) {
	world := test.NewWorld(t)

	_, err := world.ResponseWithNoBody(t.Context(), "://bad", http.MethodGet, http.Header{})
	require.Error(t, err)
}
