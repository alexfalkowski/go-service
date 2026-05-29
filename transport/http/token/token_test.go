package token_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/transport/http/token"
	"github.com/stretchr/testify/require"
)

func TestRoundTripperDoesNotMutateRequest(t *testing.T) {
	roundTripper := token.NewRoundTripper(
		env.UserID("service-user"),
		test.NewGenerator("fresh-token", nil),
		test.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			require.Equal(t, "Bearer fresh-token", req.Header.Get("Authorization"))

			return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody, Header: http.Header{}}, nil
		}),
	)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com/hello", http.NoBody)
	require.NoError(t, err)

	res, err := roundTripper.RoundTrip(req)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Empty(t, req.Header.Values("Authorization"))
}

func TestRoundTripperHandlesNilRequestHeader(t *testing.T) {
	roundTripper := token.NewRoundTripper(
		env.UserID("service-user"),
		test.NewGenerator("fresh-token", nil),
		test.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			require.Equal(t, "Bearer fresh-token", req.Header.Get("Authorization"))

			return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody, Header: http.Header{}}, nil
		}),
	)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com/hello", http.NoBody)
	require.NoError(t, err)
	req.Header = nil

	res, err := roundTripper.RoundTrip(req)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Nil(t, req.Header)
}

func TestRoundTripperClosesBodyOnGenerateError(t *testing.T) {
	roundTripper := token.NewRoundTripper(
		env.UserID("service-user"),
		test.NewGenerator("", test.ErrGenerate),
		test.RoundTripperFunc(func(*http.Request) (*http.Response, error) {
			t.Fatal("unexpected round trip")
			return nil, nil
		}),
	)
	body := &test.TrackedBody{Reader: strings.NewReader("body")}
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "http://example.com/hello", body)
	require.NoError(t, err)

	res, err := roundTripper.RoundTrip(req)
	require.Nil(t, res)
	require.Error(t, err)
	require.True(t, body.Closed)
}

func TestRoundTripperClosesBodyOnEmptyToken(t *testing.T) {
	roundTripper := token.NewRoundTripper(
		env.UserID("service-user"),
		test.NewGenerator("", nil),
		test.RoundTripperFunc(func(*http.Request) (*http.Response, error) {
			t.Fatal("unexpected round trip")
			return nil, nil
		}),
	)
	body := &test.TrackedBody{Reader: strings.NewReader("body")}
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "http://example.com/hello", body)
	require.NoError(t, err)

	res, err := roundTripper.RoundTrip(req)
	require.Nil(t, res)
	require.Error(t, err)
	require.Equal(t, http.StatusUnauthorized, status.Code(err))
	require.True(t, body.Closed)
}

func TestRoundTripperClosesBodyOnCrossOriginRedirect(t *testing.T) {
	roundTripper := token.NewRoundTripper(
		env.UserID("service-user"),
		test.NewGenerator("fresh-token", nil),
		test.RoundTripperFunc(func(*http.Request) (*http.Response, error) {
			t.Fatal("unexpected round trip")
			return nil, nil
		}),
	)
	prev, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "https://example.com/hello", http.NoBody)
	require.NoError(t, err)
	body := &test.TrackedBody{Reader: strings.NewReader("body")}
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "https://other.example.com/hello", body)
	require.NoError(t, err)
	req.Response = &http.Response{Request: prev}

	res, err := roundTripper.RoundTrip(req)
	require.Nil(t, res)
	require.ErrorIs(t, err, http.ErrUseLastResponse)
	require.True(t, body.Closed)
}
