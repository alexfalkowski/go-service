package mvc

import "github.com/alexfalkowski/go-service/v2/meta"

// Template that is rendered by the view.
type Template struct {
	Meta  meta.Map
	Model any
}
