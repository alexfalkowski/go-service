package http

import (
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

// Kind of mux for HTTP.
type Kind string

const (
	// Standard mux for HTTP.
	Standard = Kind("standard")

	// Gateway mux for HTTP.
	Gateway = Kind("gateway")
)

// NewServeMux for HTTP.
func NewServeMux(kind Kind, r *runtime.ServeMux, s *http.ServeMux) ServeMux {
	if kind == Gateway {
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
	case "Request-Id", "Geolocation", "X-Forwarded-For":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
