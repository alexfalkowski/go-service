package mvc

import (
	"github.com/alexfalkowski/go-service/v2/i18n"
	"github.com/alexfalkowski/go-service/v2/meta"
)

// Template is the composite model rendered by View templates.
type Template struct {
	Meta  meta.Map
	Model any

	// Localizer renders localized text when the server has a catalog; otherwise it is nil.
	Localizer *i18n.Localizer
}
