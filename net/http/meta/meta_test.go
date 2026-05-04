package meta_test

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/meta"
	"github.com/stretchr/testify/require"
)

func TestWithContent(t *testing.T) {
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/test", http.NoBody)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	enc := &encoder{}

	ctx := meta.WithContent(t.Context(), req, res, enc)

	require.Same(t, req, meta.Request(ctx))
	require.Same(t, res, meta.Response(ctx))
	require.Same(t, enc, meta.Encoder(ctx))
}

func TestWithContentAllowsPartialContent(t *testing.T) {
	res := httptest.NewRecorder()

	ctx := meta.WithContent(t.Context(), nil, res, nil)

	require.Same(t, res, meta.Response(ctx))
}

type encoder struct{}

func (e *encoder) Decode(_ io.Reader, _ any) error {
	return nil
}

func (e *encoder) Encode(_ io.Writer, _ any) error {
	return nil
}
