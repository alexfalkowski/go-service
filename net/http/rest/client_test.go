package rest_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/content/stream"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/net/http/rest"
	"github.com/stretchr/testify/require"
)

func TestNewClientUsesTimeout(t *testing.T) {
	rest.Register(nil, test.Content, test.StreamEncoder, test.Pool, stream.Options{})

	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set(content.TypeKey, media.Text)
		res.WriteHeader(http.StatusOK)
		res.(http.Flusher).Flush()
		<-req.Context().Done()
	}))
	t.Cleanup(server.Close)

	client := rest.NewClient(rest.WithClientTimeout("10ms"))

	err := client.Get(t.Context(), server.URL, rest.Options{})

	require.ErrorIs(t, err, context.DeadlineExceeded)
}
