package http

import (
	"net/http"

	"github.com/alexfalkowski/go-service/pkg/transport/http/breaker"
	pkgZap "github.com/alexfalkowski/go-service/pkg/transport/http/logger/zap"
	"github.com/alexfalkowski/go-service/pkg/transport/http/meta"
	"github.com/alexfalkowski/go-service/pkg/transport/http/retry"
	"github.com/alexfalkowski/go-service/pkg/transport/http/trace/opentracing"
	"go.uber.org/zap"
)

// ClientParams for HTTP.
type ClientParams struct {
	Config       *Config
	Logger       *zap.Logger
	RoundTripper http.RoundTripper
}

// NewClient for HTTP.
func NewClient(params *ClientParams) *http.Client {
	return &http.Client{Transport: newRoundTripper(params)}
}

func newRoundTripper(params *ClientParams) http.RoundTripper {
	hrt := params.RoundTripper
	if hrt == nil {
		hrt = http.DefaultTransport
	}

	hrt = pkgZap.NewRoundTripper(params.Logger, hrt)
	hrt = opentracing.NewRoundTripper(hrt)
	hrt = retry.NewRoundTripper(&params.Config.Retry, hrt)
	hrt = breaker.NewRoundTripper(hrt)
	hrt = meta.NewRoundTripper(hrt)

	return hrt
}
