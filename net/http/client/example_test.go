package client_test

import (
	"fmt"
	"net/http/httptest"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/encoding"
	encodingstream "github.com/alexfalkowski/go-service/v2/encoding/stream"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/client"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	contentstream "github.com/alexfalkowski/go-service/v2/net/http/content/stream"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-sync"
)

func ExampleClient_RequestStream() {
	pool := sync.NewBufferPool()
	streamMap := encodingstream.NewMap()
	cont := content.NewContent(encoding.NewMap(), pool)
	handler := contentstream.NewRequestHandler(
		cont,
		streamMap,
		contentstream.Options{}, func(_ context.Context, stream *contentstream.RequestStream[exampleRequest, exampleResponse]) error {
			req, err := stream.Recv()
			if err != nil {
				return err
			}

			return stream.Send(&exampleResponse{Greeting: "hello " + req.Name})
		})

	server := httptest.NewUnstartedServer(handler)
	server.EnableHTTP2 = true
	server.StartTLS()
	defer server.Close()

	c := client.NewClient(cont, streamMap, pool, client.WithRoundTripper(server.Client().Transport))
	var response exampleResponse
	err := c.RequestStream(context.Background(), http.MethodPost, server.URL, client.Options{ContentType: media.NDJSON, Accept: media.NDJSON},
		func(_ context.Context, stream *client.RequestResponseStream) error {
			if err := stream.Send(&exampleRequest{Name: "Mira"}); err != nil {
				return err
			}

			return stream.Recv(&response)
		})
	if err != nil {
		panic(err)
	}

	fmt.Println(response.Greeting)
	// Output: hello Mira
}

type exampleRequest struct {
	Name string `json:"name"`
}

type exampleResponse struct {
	Greeting string `json:"greeting"`
}
