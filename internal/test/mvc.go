package test

import (
	"embed"
	"log/slog"

	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/mvc"
)

type (
	// Todo for test.
	//nolint:godox
	Todo struct {
		Title string
		Done  bool
	}

	// Page for test.
	Page struct {
		Title string
		Todos []Todo
	}
)

var (
	// FileSystem for test.
	//go:embed views/*.tmpl
	//go:embed static/*.txt
	FileSystem embed.FS

	// Layout for test.
	Layout = mvc.NewLayout("views/full.tmpl", "views/partial.tmpl")

	// Model for test.
	Model = Page{
		Title: "My task list",
		Todos: []Todo{
			{Title: "Task 1", Done: false},
			{Title: "Task 2", Done: true},
			{Title: "Task 3", Done: true},
		},
	}
)

func registerMVC(mux *http.ServeMux, logger *slog.Logger) {
	mvc.Register(mvc.RegisterParams{
		Mux:         mux,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: logger}),
		FileSystem:  FileSystem,
		Layout:      Layout,
	})
}
