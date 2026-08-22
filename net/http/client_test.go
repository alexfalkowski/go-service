package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestClosingRoundTripperClosesBodyWhenRequested(t *testing.T) {
	t.Parallel()

	rt := http.ClosingRoundTripper(func(*http.Request) (*http.Response, error, bool) {
		return nil, io.ErrUnexpectedEOF, true
	})
	body := &test.TrackedBody{Reader: strings.NewReader("body")}
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "http://example.com", body)
	require.NoError(t, err)

	res, err := rt.RoundTrip(req)
	require.Nil(t, res)
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
	require.True(t, body.Closed)
}

func TestClosingRoundTripperLeavesDelegatedBodyOpen(t *testing.T) {
	t.Parallel()

	rt := http.ClosingRoundTripper(func(*http.Request) (*http.Response, error, bool) {
		return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody, Header: http.Header{}}, nil, false
	})
	body := &test.TrackedBody{Reader: strings.NewReader("body")}
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "http://example.com", body)
	require.NoError(t, err)

	res, err := rt.RoundTrip(req)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.False(t, body.Closed)
}
