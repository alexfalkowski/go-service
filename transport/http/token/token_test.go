package token_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/meta"
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

func TestRoundTripperGeneratesTokenForMethodPath(t *testing.T) {
	generator := &audienceGenerator{token: "fresh-token"}
	roundTripper := token.NewRoundTripper(
		env.UserID("service-user"),
		generator,
		test.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			require.Equal(t, "Bearer fresh-token", req.Header.Get("Authorization"))

			return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody, Header: http.Header{}}, nil
		}),
	)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodDelete, "http://example.com/users/123", http.NoBody)
	require.NoError(t, err)

	res, err := roundTripper.RoundTrip(req)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, "DELETE /users/123", generator.aud)
	require.Equal(t, "service-user", generator.sub)
}

func TestHandlerVerifiesTokenForMethodPath(t *testing.T) {
	verifier := &audienceVerifier{token: "fresh-token"}
	handler := token.NewHandler(env.Name("service"), env.UserID("user-id"), verifier)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPatch, "/users/123", http.NoBody)
	require.NoError(t, err)
	req = req.WithContext(meta.WithAttributes(req.Context(), meta.WithAuthorization(meta.String("fresh-token"))))
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req, func(http.ResponseWriter, *http.Request) {})

	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, "PATCH /users/123", verifier.aud)
}

func TestAccessHandler(t *testing.T) {
	for _, tt := range accessHandlerTests {
		t.Run(tt.name, func(t *testing.T) {
			handler := token.NewAccessHandler(env.Name("service"), tt.controller)
			req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, tt.path, http.NoBody)
			require.NoError(t, err)
			if tt.user != strings.Empty {
				req = req.WithContext(meta.WithAttributes(req.Context(), meta.WithUserID(meta.String(tt.user))))
			}
			res := httptest.NewRecorder()
			called := false

			handler.ServeHTTP(res, req, func(http.ResponseWriter, *http.Request) {
				called = true
			})

			require.Equal(t, tt.status, res.Code)
			require.Equal(t, tt.called, called)
		})
	}
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

type audienceGenerator struct {
	aud   string
	sub   string
	token string
}

func (g *audienceGenerator) Generate(aud, sub string) ([]byte, error) {
	g.aud = aud
	g.sub = sub

	return []byte(g.token), nil
}

type audienceVerifier struct {
	aud   string
	token string
}

func (v *audienceVerifier) Verify(token []byte, aud string) (string, error) {
	v.aud = aud
	if string(token) != v.token {
		return strings.Empty, test.ErrInvalid
	}

	return test.UserID.String(), nil
}

type accessHandlerTest struct {
	controller accessControllerFunc
	name       string
	path       string
	user       string
	status     int
	called     bool
}

var accessHandlerTests = []accessHandlerTest{
	{
		name:       "operation path bypasses access",
		path:       "/service/healthz",
		status:     http.StatusOK,
		controller: accessControllerFunc(func(context.Context) (bool, error) { return false, test.ErrInvalid }),
		called:     true,
	},
	{
		name:       "missing user id is unauthorized",
		path:       "/users/123",
		status:     http.StatusUnauthorized,
		controller: accessControllerFunc(func(context.Context) (bool, error) { return true, nil }),
	},
	{
		name:       "controller error is internal server error",
		path:       "/users/123",
		status:     http.StatusInternalServerError,
		user:       "frontend",
		controller: accessControllerFunc(func(context.Context) (bool, error) { return false, test.ErrInvalid }),
	},
	{
		name:       "access denial is forbidden",
		path:       "/users/123",
		status:     http.StatusForbidden,
		user:       "frontend",
		controller: accessControllerFunc(func(context.Context) (bool, error) { return false, nil }),
	},
	{
		name:       "access grant calls next",
		path:       "/users/123",
		status:     http.StatusOK,
		user:       "frontend",
		controller: accessControllerFunc(func(context.Context) (bool, error) { return true, nil }),
		called:     true,
	},
}

type accessControllerFunc func(context.Context) (bool, error)

func (f accessControllerFunc) HasAccess(ctx context.Context) (bool, error) {
	return f(ctx)
}
