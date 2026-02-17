package mvc

import "github.com/alexfalkowski/go-service/v2/meta"

// Template is the composite model rendered by View templates.
type Template struct {
	Meta  meta.Map
	Model any
}
