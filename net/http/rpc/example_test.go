package rpc_test

import (
	"fmt"
	"net/http/httptest"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/encoding"
	encodingstream "github.com/alexfalkowski/go-service/v2/encoding/stream"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/client"
	contentstream "github.com/alexfalkowski/go-service/v2/net/http/content/stream"
	"github.com/alexfalkowski/go-service/v2/net/http/content/unary"
	"github.com/alexfalkowski/go-service/v2/net/http/rpc"
	"github.com/alexfalkowski/go-sync"
)

func ExampleClient_Post() {
	mux := http.NewServeMux()
	router := http.NewRouter(mux, http.NewRoutePolicy())
	pool := sync.NewBufferPool()
	streamMap := encodingstream.NewMap()
	streamContent := contentstream.NewContent(streamMap, pool)
	unaryContent := unary.NewContent(encoding.NewMap(), pool)
	rpcServer := rpc.NewServer(router, unaryContent, streamContent, contentstream.Options{})

	rpcServer.Route("/hello", func(_ context.Context, req *exampleRequest) (*exampleResponse, error) {
		return &exampleResponse{Message: "hello " + req.Name}, nil
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	httpClient := client.NewClient(unaryContent, streamContent, pool)
	rpcClient := rpc.NewClient(server.URL, httpClient, rpc.WithClientContentType("application/json"))
	var res exampleResponse
	if err := rpcClient.Post(context.Background(), "/hello", &exampleRequest{Name: "Mira"}, &res); err != nil {
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
