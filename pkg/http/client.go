package http

import (
	"net/http"

	pkgZap "github.com/alexfalkowski/go-service/pkg/http/logger/zap"
	"github.com/alexfalkowski/go-service/pkg/http/meta"
	"github.com/alexfalkowski/go-service/pkg/http/trace/opentracing"
	"go.uber.org/zap"
)

// NewRoundTripper for HTTP.
func NewRoundTripper(logger *zap.Logger) http.RoundTripper {
	hrt := http.DefaultTransport
	hrt = pkgZap.NewRoundTripper(logger, hrt)
	hrt = opentracing.NewRoundTripper(hrt)
	hrt = meta.NewRoundTripper(hrt)

	return hrt
}
