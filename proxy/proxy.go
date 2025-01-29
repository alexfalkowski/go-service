package proxy

import (
	"github.com/alexfalkowski/go-service/proxy/telemetry/logger/zap"
	proxy "github.com/elazarl/goproxy"
)

// NewProxy creates a new server.
func NewProxy(logger *zap.Logger) *proxy.ProxyHttpServer {
	proxy := proxy.NewProxyHttpServer()
	proxy.AllowHTTP2 = true
	proxy.KeepDestinationHeaders = true
	proxy.KeepHeader = true
	proxy.Verbose = true
	proxy.Logger = logger

	return proxy
}
