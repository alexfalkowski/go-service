package http

import (
	"net/http"

	"go.uber.org/zap"
)

const (
	client = "client"
)

// NewRoundTripper for HTTP.
func NewRoundTripper(logger *zap.Logger) http.RoundTripper {
	hrt := http.DefaultTransport
	hrt = &loggerRoundTripper{logger: logger, RoundTripper: hrt}
	hrt = &traceRoundTripper{RoundTripper: hrt}
	hrt = &metaRoundTripper{RoundTripper: hrt}

	return hrt
}
