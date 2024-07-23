package test

import (
	"embed"
)

type (
	// Todo for test.
	//nolint:godox
	Todo struct {
		Title string
		Done  bool
	}

	// PageData for test.
	PageData struct {
		Title string
		Todos []Todo
	}
)

var (
	//go:embed views/hello.tmpl.html
	//go:embed views/error.tmpl.html
	Views embed.FS

	// Model for test.
	Model = PageData{
		Title: "My TODO list",
		Todos: []Todo{
			{Title: "Task 1", Done: false},
			{Title: "Task 2", Done: true},
			{Title: "Task 3", Done: true},
		},
	}
)
