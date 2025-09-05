package test

import (
	"github.com/alexfalkowski/go-health/v2/checker"
	"github.com/alexfalkowski/go-health/v2/server"
	"github.com/alexfalkowski/go-service/v2/cache/cacher"
	"github.com/alexfalkowski/go-service/v2/config"
	"github.com/alexfalkowski/go-service/v2/crypto/aes"
	"github.com/alexfalkowski/go-service/v2/crypto/bcrypt"
	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/crypto/hmac"
	"github.com/alexfalkowski/go-service/v2/crypto/rsa"
	"github.com/alexfalkowski/go-service/v2/crypto/ssh"
	"github.com/alexfalkowski/go-service/v2/database/sql/pg"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/feature"
	"github.com/alexfalkowski/go-service/v2/health"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/module"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/time"
	gt "github.com/alexfalkowski/go-service/v2/transport/grpc/token"
	"github.com/alexfalkowski/go-service/v2/transport/http"
	"github.com/alexfalkowski/go-service/v2/transport/http/events"
	ht "github.com/alexfalkowski/go-service/v2/transport/http/token"
	"github.com/open-feature/go-sdk/openfeature"
	webhooks "github.com/standard-webhooks/standard-webhooks/libraries/go"
	"google.golang.org/grpc"
)

// Options for test.
func Options() []di.Option {
	return []di.Option{
		module.Server,
		di.Decorate(decorateConfig),
		di.Constructor(registrations),
		di.Register(healthRegister), di.Register(healthObserver),
		di.Register(livenessObserver), di.Register(readinessObserver),
		di.Register(grpcObserver), di.Register(invokeServiceRegistrar),
		di.Register(shutdown), di.Register(invokeFeatureClient), di.Register(invokeWebhooks), di.Register(invokeConfigs),
		di.Register(invokeMeter), di.Register(invokeNetwork), di.Register(invokeCache),
		di.Register(invokeCrypt), di.Register(invokeEnvironment), di.Register(invokeTokens),
		di.Register(invokeAccessController),
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

func healthRegister(name env.Name, server *server.Server, regs health.Registrations) {
	server.Register(name.String(), regs...)
}

func healthObserver(name env.Name, server *server.Server) error {
	return server.Observe(name.String(), "healthz", "noop")
}

func livenessObserver(name env.Name, server *server.Server) error {
	return server.Observe(name.String(), "livez", "noop")
}

func readinessObserver(name env.Name, server *server.Server) error {
	return server.Observe(name.String(), "readyz", "http", "online")
}

func grpcObserver(name env.Name, server *server.Server) error {
	return server.Observe(name.String(), "grpc", "http")
}

func decorateConfig(cfg *config.Config) *config.Config {
	return cfg
}

func invokeServiceRegistrar(_ grpc.ServiceRegistrar) {}

func invokeCache(_ cacher.Cache) {}

func invokeConfigs(_ *pg.Config, _ *feature.Config, _ *id.Config) {}

func invokeMeter(_ *metrics.Meter) {}

func invokeFeatureClient(_ *openfeature.Client) {}

func invokeWebhooks(_ *webhooks.Webhook, _ *events.Receiver) {}

func invokeEnvironment(_ env.Name, _ env.UserAgent, _ env.Version) {}

func invokeNetwork(_ time.Network) {}

func invokeAccessController(_ ht.AccessController, _ gt.AccessController) {}

func invokeTokens(_ ht.Generator, _ ht.Verifier, _ gt.Generator, _ gt.Verifier) {}

func invokeCrypt(
	signer *bcrypt.Signer,
	_ *ed25519.Signer, _ *ed25519.Verifier,
	_ *rsa.Encryptor, _ *rsa.Decryptor,
	_ *aes.Cipher,
	_ *hmac.Signer,
	_ *ssh.Signer, _ *ssh.Verifier,
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

func shutdown(s di.Shutdowner) {
	go func(s di.Shutdowner) {
		time.Sleep(time.Second)

		_ = s.Shutdown()
	}(s)
}
