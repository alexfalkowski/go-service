package http

import (
	"net/http"

	pkgZap "github.com/alexfalkowski/go-service/pkg/http/logger/zap"
	"go.uber.org/zap"
)

// NewRoundTripper for HTTP.
func NewRoundTripper(logger *zap.Logger) http.RoundTripper {
	hrt := http.DefaultTransport
	hrt = pkgZap.NewRoundTripper(logger, hrt)
	hrt = &traceRoundTripper{RoundTripper: hrt}
	hrt = &metaRoundTripper{RoundTripper: hrt}

	return hrt
}
