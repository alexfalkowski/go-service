package test

import (
	"time"

	"github.com/alexfalkowski/go-health/checker"
	"github.com/alexfalkowski/go-health/server"
	"github.com/alexfalkowski/go-service/cache"
	"github.com/alexfalkowski/go-service/cache/cacher"
	"github.com/alexfalkowski/go-service/cli"
	"github.com/alexfalkowski/go-service/config"
	"github.com/alexfalkowski/go-service/crypto/aes"
	"github.com/alexfalkowski/go-service/crypto/bcrypt"
	"github.com/alexfalkowski/go-service/crypto/ed25519"
	"github.com/alexfalkowski/go-service/crypto/hmac"
	"github.com/alexfalkowski/go-service/crypto/rsa"
	cs "github.com/alexfalkowski/go-service/crypto/ssh"
	"github.com/alexfalkowski/go-service/database/sql"
	"github.com/alexfalkowski/go-service/database/sql/pg"
	sd "github.com/alexfalkowski/go-service/debug"
	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/feature"
	"github.com/alexfalkowski/go-service/health"
	shg "github.com/alexfalkowski/go-service/health/transport/grpc"
	shh "github.com/alexfalkowski/go-service/health/transport/http"
	"github.com/alexfalkowski/go-service/hooks"
	"github.com/alexfalkowski/go-service/id"
	"github.com/alexfalkowski/go-service/module"
	"github.com/alexfalkowski/go-service/strings"
	"github.com/alexfalkowski/go-service/telemetry"
	"github.com/alexfalkowski/go-service/telemetry/logger"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	st "github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/token"
	"github.com/alexfalkowski/go-service/token/jwt"
	"github.com/alexfalkowski/go-service/token/paseto"
	ts "github.com/alexfalkowski/go-service/token/ssh"
	"github.com/alexfalkowski/go-service/transport"
	geh "github.com/alexfalkowski/go-service/transport/events/http"
	"github.com/alexfalkowski/go-service/transport/http"
	"github.com/open-feature/go-sdk/openfeature"
	h "github.com/standard-webhooks/standard-webhooks/libraries/go"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

// Options for test.
func Options() []fx.Option {
	return []fx.Option{
		module.Module, cli.Module, config.Module, sd.Module,
		feature.Module, transport.Module, telemetry.Module, health.Module,
		sql.Module, hooks.Module, cache.Module, token.Module,
		fx.Provide(registrations), fx.Provide(healthObserver), fx.Provide(livenessObserver),
		fx.Provide(readinessObserver), fx.Provide(grpcObserver), fx.Invoke(invokeServiceRegistrar),
		fx.Invoke(shutdown), fx.Invoke(invokeFeatureClient), fx.Invoke(invokeWebhooks), fx.Invoke(invokeConfigs),
		fx.Invoke(invokeMeter), fx.Invoke(invokeNetwork), fx.Invoke(invokeCache),
		fx.Invoke(invokeCrypt), fx.Invoke(invokeEnvironment), fx.Invoke(invokeTokens),
	}
}

func registrations(logger *logger.Logger, cfg *http.Config, ua env.UserAgent, tracer *tracer.Tracer, _ env.Version) health.Registrations {
	if cfg == nil {
		return nil
	}

	timeout := 5 * time.Second
	nc := checker.NewNoopChecker()
	nr := server.NewRegistration("noop", timeout, nc)
	rt, _ := http.NewRoundTripper(http.WithClientLogger(logger), http.WithClientTracer(tracer), http.WithClientUserAgent(ua))
	hc := checker.NewHTTPChecker("https://google.com", timeout, checker.WithRoundTripper(rt))
	hr := server.NewRegistration("http", timeout, hc)

	return health.Registrations{nr, hr, server.NewOnlineRegistration(timeout, timeout)}
}

func healthObserver(healthServer *server.Server) (*shh.HealthObserver, error) {
	return &shh.HealthObserver{Observer: healthServer.Observe("noop")}, nil
}

func livenessObserver(healthServer *server.Server) *shh.LivenessObserver {
	return &shh.LivenessObserver{Observer: healthServer.Observe("noop")}
}

func readinessObserver(healthServer *server.Server) *shh.ReadinessObserver {
	return &shh.ReadinessObserver{Observer: healthServer.Observe("http", "online")}
}

func grpcObserver(healthServer *server.Server) *shg.Observer {
	return &shg.Observer{Observer: healthServer.Observe("http")}
}

func invokeServiceRegistrar(_ grpc.ServiceRegistrar) {}

func invokeCache(_ cacher.Cache) {}

func invokeConfigs(_ *pg.Config, _ *feature.Config, _ *id.Config) {}

func invokeMeter(_ *metrics.Meter) {}

func invokeFeatureClient(_ *openfeature.Client) {}

func invokeWebhooks(_ *h.Webhook, _ *geh.Receiver) {}

func invokeEnvironment(_ env.Name, _ env.UserAgent, _ env.Version) {}

func invokeNetwork(_ st.Network) {}

func invokeCrypt(
	signer *bcrypt.Signer,
	_ *ed25519.Signer, _ *ed25519.Verifier,
	_ *rsa.Encryptor, _ *rsa.Decryptor,
	_ *aes.Cipher,
	_ *hmac.Signer,
	_ *cs.Signer, _ *cs.Verifier,
) error {
	msg := strings.Bytes("hello")

	e, err := signer.Sign(msg)
	if err != nil {
		return err
	}

	err = signer.Verify(e, msg)
	if err != nil {
		return err
	}

	return nil
}

func invokeTokens(_ *jwt.Token, _ *paseto.Token, _ *ts.Token, _ *token.Token) {}

func shutdown(s fx.Shutdowner) {
	go func(s fx.Shutdowner) {
		time.Sleep(time.Second)

		_ = s.Shutdown()
	}(s)
}
