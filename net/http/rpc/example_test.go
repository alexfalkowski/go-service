package rpc_test

import (
	"fmt"
	"net/http/httptest"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/encoding/json"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/rpc"
	"github.com/alexfalkowski/go-sync"
)

func ExampleClient_Post() {
	mux := http.NewServeMux()
	pool := sync.NewBufferPool()
	rpc.Register(rpc.RegisterParams{
		Mux:     mux,
		Content: exampleContent(pool),
		Pool:    pool,
	})

	rpc.Route("/hello", func(_ context.Context, req *exampleRequest) (*exampleResponse, error) {
		return &exampleResponse{Message: "hello " + req.Name}, nil
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := rpc.NewClient(server.URL, rpc.WithClientContentType("application/json"))
	var res exampleResponse
	if err := client.Post(context.Background(), "/hello", &exampleRequest{Name: "Mira"}, &res); err != nil {
		panic(err)
	}

	fmt.Println(res.Message)
	// Output: hello Mira
}

type exampleRequest struct {
	Name string `json:"name"`
}

type exampleResponse struct {
	Message string `json:"message"`
}

func exampleContent(pool *sync.BufferPool) *content.Content {
	return content.NewContent(
		encoding.NewMap(encoding.MapParams{JSON: json.NewEncoder()}),
		pool,
	)
}
