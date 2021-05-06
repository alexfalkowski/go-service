package http

import (
	"net/http"

	pkgZap "github.com/alexfalkowski/go-service/pkg/transport/http/logger/zap"
	"github.com/alexfalkowski/go-service/pkg/transport/http/meta"
	"github.com/alexfalkowski/go-service/pkg/transport/http/trace/opentracing"
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
