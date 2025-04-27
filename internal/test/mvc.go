package test

import "embed"

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
	//go:embed views/*.tmpl
	//go:embed static/*.txt
	Views embed.FS

	// Model for test.
	Model = Page{
		Title: "My TODO list",
		Todos: []Todo{
			{Title: "Task 1", Done: false},
			{Title: "Task 2", Done: true},
			{Title: "Task 3", Done: true},
		},
	}
)
