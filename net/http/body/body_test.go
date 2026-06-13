package body_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/body"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestReadAllHandlesNilBody(t *testing.T) {
	req := &http.Request{}

	data, bufferedBody, err := body.ReadAll(req)
	require.NoError(t, err)
	require.Empty(t, data)
	require.NotNil(t, bufferedBody)
	require.Equal(t, http.NoBody, req.Body)
}

func TestReadAllBuffersBody(t *testing.T) {
	req := &http.Request{Body: &test.TrackedBody{Reader: strings.NewReader("body")}}

	data, bufferedBody, err := body.ReadAll(req)
	require.NoError(t, err)
	require.Equal(t, []byte("body"), data)
	require.NotNil(t, bufferedBody)
}

func TestCloseSkipsEmptyBody(t *testing.T) {
	body.Close(nil)
	body.Close(http.NoBody)
}

func TestCloseClosesBody(t *testing.T) {
	trackedBody := &test.TrackedBody{Reader: strings.NewReader("body")}

	body.Close(trackedBody)
	require.True(t, trackedBody.Closed)
}
