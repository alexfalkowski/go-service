package rest_test

import (
	"fmt"
	"net/http/httptest"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/encoding"
	encodingstream "github.com/alexfalkowski/go-service/v2/encoding/stream"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	contentstream "github.com/alexfalkowski/go-service/v2/net/http/content/stream"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/net/http/rest"
	"github.com/alexfalkowski/go-sync"
)

func ExampleGet() {
	mux := http.NewServeMux()
	router := http.NewRouter(mux, http.NewRoutePolicy())
	pool := sync.NewBufferPool()
	streamMap := encodingstream.NewMap()
	rest.Register(router, content.NewContent(encoding.NewMap(), pool), streamMap, pool, contentstream.Options{})

	rest.Get("/hello", func(context.Context) (*exampleResponse, error) {
		return &exampleResponse{Message: "hello"}, nil
	})

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/hello", http.NoBody)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	fmt.Println(res.Code)
	fmt.Print(res.Body.String())
	// Output:
	// 200
	// {
	//   "message": "hello"
	// }
}

func ExampleStreamGet() {
	mux := http.NewServeMux()
	router := http.NewRouter(mux, http.NewRoutePolicy())
	pool := sync.NewBufferPool()
	streamMap := encodingstream.NewMap()
	rest.Register(router, content.NewContent(encoding.NewMap(), pool), streamMap, pool, contentstream.Options{})

	rest.StreamGet("/hello", func(_ context.Context, stream *contentstream.Stream[exampleResponse]) error {
		if err := stream.Send(&exampleResponse{Message: "hello"}); err != nil {
			return err
		}

		return stream.Send(&exampleResponse{Message: "world"})
	})

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	fmt.Println(res.Code)
	fmt.Print(res.Body.String())
	// Output:
	// 200
	// {"message":"hello"}
	// {"message":"world"}
}

type exampleResponse struct {
	Message string `json:"message"`
}
