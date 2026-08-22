package rest_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/client"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/net/http/rest"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestNewClientUsesTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set(http.ContentTypeKey, media.Text)
		res.WriteHeader(http.StatusOK)
		res.(http.Flusher).Flush()
		<-req.Context().Done()
	}))
	t.Cleanup(server.Close)

	restClient := rest.NewClient(test.NewContentClient(client.WithTimeout(time.MustParseDuration("10ms"))))

	err := restClient.Get(t.Context(), server.URL, rest.Options{})

	require.ErrorIs(t, err, context.DeadlineExceeded)
}
