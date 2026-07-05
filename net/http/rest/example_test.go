package rest_test

import (
	"fmt"
	"net/http/httptest"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/encoding/json"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/rest"
	"github.com/alexfalkowski/go-sync"
)

func ExampleGet() {
	mux := http.NewServeMux()
	pool := sync.NewBufferPool()
	rest.Register(rest.RegisterParams{
		Mux: mux,
		Content: content.NewContent(
			encoding.NewMap(encoding.MapParams{JSON: json.NewEncoder()}),
			pool,
		),
		Pool: pool,
	})

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

type exampleResponse struct {
	Message string `json:"message"`
}
