package mvc

import (
	"github.com/alexfalkowski/go-service/v3/context"
	"github.com/alexfalkowski/go-service/v3/i18n"
	"github.com/alexfalkowski/go-service/v3/net/http"
	"github.com/alexfalkowski/go-service/v3/net/http/meta"
	"github.com/alexfalkowski/go-service/v3/strings"
)

const languageKey = context.Key("mvc-language")

// WithLanguage returns a context whose language preference takes precedence over
// Accept-Language. Set it in middleware before MVC handlers run, or pass the
// returned context to View.Render. An empty, invalid, or unsupported preference is ignored.
// Applications own URL, cookie, and user-profile selection policies.
func WithLanguage(ctx context.Context, language string) context.Context {
	return context.WithValue(ctx, languageKey, language)
}

func (v *View) localizer(ctx context.Context) *i18n.Localizer {
	if v.catalog == nil {
		return nil
	}

	preference, _ := ctx.Value(languageKey).(string)
	preferences := []string{preference}
	if req := meta.Request(ctx); req != nil {
		preferences = append(preferences, strings.Join(",", req.Header.Values(http.AcceptLanguageKey)...))
	}

	localizer := v.catalog.Localizer(preferences...)
	header := meta.Response(ctx).Header()
	http.AddVary(header, http.AcceptLanguageKey)
	header.Set(http.ContentLanguageKey, localizer.Tag().String())

	return localizer
}
