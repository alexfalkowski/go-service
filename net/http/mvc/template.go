package mvc

import "github.com/alexfalkowski/go-service/meta"

// Template that is rendered by the view.
type Template struct {
	Meta  meta.Map
	Model any
}
