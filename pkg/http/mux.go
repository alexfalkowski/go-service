package http

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// NewMux for HTTP.
func NewMux() *runtime.ServeMux {
	return runtime.NewServeMux()
}
