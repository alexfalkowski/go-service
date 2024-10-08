package test

import (
	"net/http"

	lm "github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	v1 "github.com/alexfalkowski/go-service/test/greet/v1"
	"github.com/alexfalkowski/go-service/token"
	"github.com/alexfalkowski/go-service/transport"
	tg "github.com/alexfalkowski/go-service/transport/grpc"
	th "github.com/alexfalkowski/go-service/transport/http"
	"github.com/sethvargo/go-limiter"
	"github.com/urfave/negroni/v3"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Server for test.
type Server struct {
	Lifecycle  fx.Lifecycle
	Meter      metric.Meter
	Verifier   token.Verifier
	Mux        *http.ServeMux
	HTTP       *th.Server
	GRPC       *tg.Server
	Transport  *transport.Config
	Tracer     *tracer.Config
	Limiter    limiter.Store
	Key        lm.KeyFunc
	Logger     *zap.Logger
	VerifyAuth bool
}

// Register server.
func (s *Server) Register() {
	sh := NewShutdowner()
	tracer, err := tracer.NewTracer(s.Lifecycle, Environment, Version, Name, s.Tracer, s.Logger)
	runtime.Must(err)

	h, err := th.NewServer(th.ServerParams{
		Shutdowner: sh, Mux: s.Mux,
		Config: s.Transport.HTTP, Logger: s.Logger,
		Tracer: tracer, Meter: s.Meter,
		Limiter: s.Limiter, Key: s.Key, Handlers: []negroni.Handler{&none{}},
		Verifier: s.Verifier, UserAgent: UserAgent, Version: Version,
	})
	runtime.Must(err)

	s.HTTP = h

	g, err := tg.NewServer(tg.ServerParams{
		Shutdowner: sh, Config: s.Transport.GRPC, Logger: s.Logger,
		Tracer: tracer, Meter: s.Meter, Limiter: s.Limiter, Key: s.Key,
		Verifier: s.Verifier, UserAgent: UserAgent, Version: Version,
	})
	runtime.Must(err)

	s.GRPC = g

	v1.RegisterGreeterServiceServer(g.Server(), NewService(s.VerifyAuth))
	transport.Register(transport.RegisterParams{Lifecycle: s.Lifecycle, Servers: []transport.Server{h, g}})
}

type none struct{}

func (*none) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(rw, r)
}
