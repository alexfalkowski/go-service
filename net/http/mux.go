package http

import (
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

// MuxKind for HTTP.
type MuxKind string

const (
	// StandardMux for HTTP.
	StandardMux = MuxKind("standard")

	// GatewayMux for HTTP.
	GatewayMux = MuxKind("gateway")
)

// NewServeMux for HTTP.
func NewServeMux(kind MuxKind, r *runtime.ServeMux, s *http.ServeMux) ServeMux {
	if kind == GatewayMux {
		return &RuntimeServeMux{r}
	}

	return &StandardServeMux{s}
}

// ServeMux for HTTP.
type ServeMux interface {
	// Handle a verb, pattern with the func.
	Handle(verb, pattern string, fn http.HandlerFunc) error

	// Handler from the mux.
	Handler() http.Handler
}

// StandardServeMux for HTTP.
type StandardServeMux struct {
	*http.ServeMux
}

func (s *StandardServeMux) Handle(verb, pattern string, fn http.HandlerFunc) error {
	s.HandleFunc(fmt.Sprintf("%s %s", verb, pattern), fn)

	return nil
}

func (s *StandardServeMux) Handler() http.Handler {
	return s.ServeMux
}

// RuntimeServeMux for HTTP.
type RuntimeServeMux struct {
	*runtime.ServeMux
}

func (r *RuntimeServeMux) Handle(verb, pattern string, fn http.HandlerFunc) error {
	return r.HandlePath(verb, pattern, func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		fn(w, r)
	})
}

func (r *RuntimeServeMux) Handler() http.Handler {
	return r.ServeMux
}

// NewServeMux for http.
func NewStandardServeMux() *http.ServeMux {
	return http.NewServeMux()
}

// NewRuntimeServeMux for HTTP.
func NewRuntimeServeMux() *runtime.ServeMux {
	opts := []runtime.ServeMuxOption{
		runtime.WithIncomingHeaderMatcher(customMatcher),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	}

	return runtime.NewServeMux(opts...)
}

func customMatcher(key string) (string, bool) {
	switch key {
	case "X-Real-IP", "CF-Connecting-IP", "True-Client-IP", "X-Forwarded-For":
		return key, true
	case "Request-Id", "Geolocation":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
