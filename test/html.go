package test

import (
	"embed"
)

type (
	Todo struct {
		Title string
		Done  bool
	}

	PageData struct {
		Title string
		Todos []Todo
	}
)

var (
	//go:embed html/hello.tmpl
	//go:embed html/error.tmpl
	HTML embed.FS

	// HTMLData for test.
	HTMLData = PageData{
		Title: "My TODO list",
		Todos: []Todo{
			{Title: "Task 1", Done: false},
			{Title: "Task 2", Done: true},
			{Title: "Task 3", Done: true},
		},
	}
)
