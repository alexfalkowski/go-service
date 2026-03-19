package test

import (
	"embed"
	"log/slog"

	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/mvc"
)

type (
	// Todo is the view model for a single todo item in the embedded MVC fixtures.
	//nolint:godox
	Todo struct {
		Title string
		Done  bool
	}

	// Page is the top-level view model rendered by the embedded MVC fixtures.
	Page struct {
		Title string
		Todos []Todo
	}
)

var (
	// FileSystem embeds the MVC templates and static assets used by tests.
	//go:embed views/*.tmpl
	//go:embed static/*.txt
	FileSystem embed.FS

	// Layout is the shared MVC layout used when rendering the embedded templates.
	Layout = mvc.NewLayout("views/full.tmpl", "views/partial.tmpl")

	// Model is the sample page rendered by MVC tests.
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
