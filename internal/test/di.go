package test

import (
	"github.com/alexfalkowski/go-health/v2/checker"
	"github.com/alexfalkowski/go-health/v2/server"
	"github.com/alexfalkowski/go-service/v3/cache"
	"github.com/alexfalkowski/go-service/v3/config"
	"github.com/alexfalkowski/go-service/v3/crypto/aes"
	"github.com/alexfalkowski/go-service/v3/crypto/bcrypt"
	"github.com/alexfalkowski/go-service/v3/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v3/crypto/hmac"
	"github.com/alexfalkowski/go-service/v3/crypto/rsa"
	"github.com/alexfalkowski/go-service/v3/database/sql"
	"github.com/alexfalkowski/go-service/v3/database/sql/pg"
	"github.com/alexfalkowski/go-service/v3/env"
	"github.com/alexfalkowski/go-service/v3/feature"
	"github.com/alexfalkowski/go-service/v3/health"
	healthchecker "github.com/alexfalkowski/go-service/v3/health/checker"
	"github.com/alexfalkowski/go-service/v3/hooks"
	"github.com/alexfalkowski/go-service/v3/id"
	"github.com/alexfalkowski/go-service/v3/module"
	"github.com/alexfalkowski/go-service/v3/net/grpc"
	"github.com/alexfalkowski/go-service/v3/strings"
	"github.com/alexfalkowski/go-service/v3/telemetry/logger"
	"github.com/alexfalkowski/go-service/v3/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v3/time"
	"github.com/alexfalkowski/go-service/v3/token/access"
	grpctoken "github.com/alexfalkowski/go-service/v3/transport/grpc/token"
	"github.com/alexfalkowski/go-service/v3/transport/http"
	"github.com/alexfalkowski/go-service/v3/transport/http/events"
	httptoken "github.com/alexfalkowski/go-service/v3/transport/http/token"
	"github.com/alexfalkowski/go-signal/v2"
	"github.com/open-feature/go-sdk/openfeature"
)

// NoopLifecycle is a [signal.Lifecycle] test double that ignores appended hooks.
type NoopLifecycle struct{}

// Append implements [signal.Lifecycle].
func (NoopLifecycle) Append(signal.Hook) {}

// Options returns the DI options used by tests that exercise the full server module wiring.
//
// The option set decorates config, registers health checks and observers, and
// forces construction of the major infrastructure components so go-signal tests
// can verify end-to-end module composition.
func Options() []signal.Option {
	return []signal.Option{
		module.Server,
		signal.Decorate(decorateConfig),
		signal.Constructor(registrations),
		signal.Register(healthRegister), signal.Register(healthObserver),
		signal.Register(livenessObserver), signal.Register(readinessObserver),
		signal.Register(grpcObserver), signal.Register(invokeServiceRegistrar),
		signal.Register(shutdown), signal.Register(invokeFeatureClient), signal.Register(invokeWebhooks), signal.Register(invokeConfigs),
		signal.Register(invokeMeter), signal.Register(invokeNetwork), signal.Register(invokeCache),
		signal.Register(invokeCrypt), signal.Register(invokeEnvironment), signal.Register(invokeTokens),
		signal.Register(invokeAccessController), signal.Register(invokeDB),
	}
}

func registrations(logger *logger.Logger, cfg *http.Config, pinger cache.Pinger, ua env.UserAgent, _ env.Version) (health.Registrations, error) {
	if cfg == nil {
		return nil, nil
	}

	timeout := 5 * time.Second
	nc := checker.NewNoopChecker()
	nr := server.NewRegistration("noop", timeout.Duration(), nc)
	rt, err := http.NewRoundTripper(http.WithClientLogger(logger), http.WithClientUserAgent(ua))
	if err != nil {
		return nil, err
	}

	hc := checker.NewHTTPChecker(StatusURL("200"), timeout.Duration(), checker.WithRoundTripper(rt))
	hr := server.NewRegistration("http", timeout.Duration(), hc)

	regs := health.Registrations{nr, hr, server.NewOnlineRegistration(timeout.Duration(), timeout.Duration())}
	if pinger != nil {
		cc := healthchecker.NewCacheChecker(pinger, timeout)
		regs = append(regs, server.NewRegistration("cache", timeout.Duration(), cc))
	}

	return regs, nil
}

func healthRegister(cfg *http.Config, name env.Name, server *server.Server, regs health.Registrations) {
	if cfg == nil {
		return
	}

	server.Register(name.String(), regs...)
}

func healthObserver(cfg *http.Config, name env.Name, server *server.Server) error {
	if cfg == nil {
		return nil
	}

	return server.Observe(name.String(), "healthz", "noop")
}

func livenessObserver(cfg *http.Config, name env.Name, server *server.Server) error {
	if cfg == nil {
		return nil
	}

	return server.Observe(name.String(), "livez", "noop")
}

func readinessObserver(cfg *http.Config, pinger cache.Pinger, name env.Name, server *server.Server) error {
	if cfg == nil {
		return nil
	}

	if pinger != nil {
		return server.Observe(name.String(), "readyz", "http", "online", "cache")
	}

	return server.Observe(name.String(), "readyz", "http", "online")
}

func grpcObserver(cfg *http.Config, name env.Name, server *server.Server) error {
	if cfg == nil {
		return nil
	}

	return server.Observe(name.String(), "grpc", "http")
}

func decorateConfig(cfg *config.Config) *config.Config {
	return cfg
}

func invokeServiceRegistrar(_ grpc.ServiceRegistrar) {}

func invokeCache(_ *cache.Cache) {}

func invokeConfigs(_ *pg.Config, _ *feature.Config, _ *id.Config) {}

func invokeMeter(_ metrics.Meter) {}

func invokeFeatureClient(_ *openfeature.Client) {}

func invokeWebhooks(_ *hooks.Hook, _ *events.Receiver) {}

func invokeEnvironment(_ env.Name, _ env.UserAgent, _ env.Version) {}

func invokeNetwork(_ time.Network) {}

func invokeAccessController(_ access.Controller) {}

func invokeTokens(_ httptoken.Generator, _ httptoken.Verifier, _ grpctoken.Generator, _ grpctoken.Verifier) {
}

func invokeDB(_ *sql.DBs) {}

func invokeCrypt(
	signer *bcrypt.Signer,
	_ *ed25519.Signer, _ *ed25519.Verifier,
	_ *rsa.Encryptor, _ *rsa.Decryptor,
	_ *aes.Cipher,
	_ *hmac.Signer,
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

func shutdown(s signal.Shutdowner) {
	go func(s signal.Shutdowner) {
		time.Sleep(time.Second)

		s.Shutdown(0)
	}(s)
}
