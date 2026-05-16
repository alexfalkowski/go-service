package mvc_test

import (
	"fmt"
	"log/slog"
	"net/http/httptest"
	"testing/fstest"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/mvc"
	sync "github.com/alexfalkowski/go-sync"
)

func ExampleGet() {
	mux := http.NewServeMux()
	mvc.Register(mvc.RegisterParams{
		Mux:         mux,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem:  exampleFileSystem(),
		Pool:        sync.NewBufferPool(),
		Layout:      mvc.NewLayout("views/full.tmpl", "views/partial.tmpl"),
	})

	mvc.Get("/hello", func(_ context.Context) (*mvc.View, *examplePage, error) {
		return mvc.NewFullView("views/page.tmpl"), &examplePage{Title: "Hello"}, nil
	})

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/hello", http.NoBody)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	fmt.Println(res.Code)
	fmt.Println(res.Body.String())
	// Output:
	// 200
	// Hello
}

func ExampleNotFound() {
	mux := http.NewServeMux()
	mvc.Register(mvc.RegisterParams{
		Mux:         mux,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem:  exampleFileSystem(),
		Pool:        sync.NewBufferPool(),
		Layout:      mvc.NewLayout("views/full.tmpl", "views/partial.tmpl"),
	})

	mvc.NotFound(func(_ context.Context) (*mvc.View, *examplePage) {
		return mvc.NewFullView("views/page.tmpl"), &examplePage{Title: "Not Found"}
	})

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/missing", http.NoBody)
	res := httptest.NewRecorder()

	mvc.NewHandler(mux).ServeHTTP(res, req)

	fmt.Println(res.Code)
	fmt.Println(res.Body.String())
	// Output:
	// 404
	// Not Found
}

type examplePage struct {
	Title string
}

func exampleFileSystem() fstest.MapFS {
	return fstest.MapFS{
		"views/full.tmpl":    &fstest.MapFile{Data: []byte(`{{ block "content" . }}{{ end }}`)},
		"views/partial.tmpl": &fstest.MapFile{Data: []byte(`{{ block "content" . }}{{ end }}`)},
		"views/page.tmpl":    &fstest.MapFile{Data: []byte(`{{ define "content" }}{{ .Model.Title }}{{ end }}`)},
	}
}
