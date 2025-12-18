package hooks_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/transport/http/hooks"
	"github.com/stretchr/testify/require"
)

func TestVerify(t *testing.T) {
	hook := hooks.NewWebhook(nil, nil)
	req := &http.Request{Body: &test.ErrReaderCloser{}}

	require.Error(t, hook.Verify(req))
}

func TestSign(t *testing.T) {
	hook := hooks.NewWebhook(nil, nil)
	req := &http.Request{Body: &test.ErrReaderCloser{}}

	require.Error(t, hook.Sign(req))
}

func TestRoundTripper(t *testing.T) {
	hook := hooks.NewWebhook(nil, nil)
	rt := hooks.NewRoundTripper(hook, nil)
	req := &http.Request{Body: &test.ErrReaderCloser{}}

	_, err := rt.RoundTrip(req)
	require.Error(t, err)
}
