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
	router := http.NewRouter(mux, http.NewRoutePolicy())
	mvc.Register(mvc.RegisterParams{
		Router:      router,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem:  exampleFileSystem(),
		Pool:        sync.NewBufferPool(),
		Layout:      mvc.NewLayout("views/full.tmpl", "views/partial.tmpl"),
	})

	view := mvc.NewFullView("views/page.tmpl")
	mvc.Get("/hello", func(_ context.Context) (*mvc.View, *examplePage, error) {
		return view, &examplePage{Title: "Hello"}, nil
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
	router := http.NewRouter(mux, http.NewRoutePolicy())
	mvc.Register(mvc.RegisterParams{
		Router:      router,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem:  exampleFileSystem(),
		Pool:        sync.NewBufferPool(),
		Layout:      mvc.NewLayout("views/full.tmpl", "views/partial.tmpl"),
	})

	view := mvc.NewFullView("views/page.tmpl")
	mvc.NotFound(func(_ context.Context) (*mvc.View, *examplePage) {
		return view, &examplePage{Title: "Not Found"}
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

func ExampleStaticFile() {
	mux := http.NewServeMux()
	router := http.NewRouter(mux, http.NewRoutePolicy())
	mvc.Register(mvc.RegisterParams{
		Router:      router,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem:  exampleFileSystem(),
		Pool:        sync.NewBufferPool(),
		Layout:      mvc.NewLayout("views/full.tmpl", "views/partial.tmpl"),
	})

	mvc.StaticFile(
		"/asset.txt",
		"static/asset.txt",
		mvc.WithCacheControl("public, max-age=60"),
	)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/asset.txt", http.NoBody)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	fmt.Println(res.Code)
	fmt.Println(res.Header().Get("Cache-Control"))
	fmt.Println(res.Body.String())
	// Output:
	// 200
	// public, max-age=60
	// hello
}

func ExampleStaticPathValue() {
	mux := http.NewServeMux()
	router := http.NewRouter(mux, http.NewRoutePolicy())
	mvc.Register(mvc.RegisterParams{
		Router:      router,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem:  exampleFileSystem(),
		Pool:        sync.NewBufferPool(),
		Layout:      mvc.NewLayout("views/full.tmpl", "views/partial.tmpl"),
	})

	mvc.StaticPathValue("/{file...}", "file", "static")

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/asset.txt", http.NoBody)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	fmt.Println(res.Code)
	fmt.Println(res.Body.String())
	// Output:
	// 200
	// hello
}

type examplePage struct {
	Title string
}

func exampleFileSystem() fstest.MapFS {
	return fstest.MapFS{
		"views/full.tmpl":    &fstest.MapFile{Data: []byte(`{{ block "content" . }}{{ end }}`)},
		"views/partial.tmpl": &fstest.MapFile{Data: []byte(`{{ block "content" . }}{{ end }}`)},
		"views/page.tmpl":    &fstest.MapFile{Data: []byte(`{{ define "content" }}{{ .Model.Title }}{{ end }}`)},
		"static/asset.txt":   &fstest.MapFile{Data: []byte("hello")},
	}
}
