package http

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// NewMux for HTTP.
func NewMux() *runtime.ServeMux {
	return runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(customMatcher))
}

func customMatcher(key string) (string, bool) {
	switch key {
	case "Request-Id":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
