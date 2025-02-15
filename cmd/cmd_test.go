package cmd_test

import (
	"testing"
	"time"

	"github.com/alexfalkowski/go-health/checker"
	"github.com/alexfalkowski/go-health/server"
	"github.com/alexfalkowski/go-service/cache"
	cc "github.com/alexfalkowski/go-service/cache/config"
	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/config"
	"github.com/alexfalkowski/go-service/crypto/aes"
	"github.com/alexfalkowski/go-service/crypto/argon2"
	"github.com/alexfalkowski/go-service/crypto/ed25519"
	"github.com/alexfalkowski/go-service/crypto/hmac"
	"github.com/alexfalkowski/go-service/crypto/rsa"
	"github.com/alexfalkowski/go-service/crypto/ssh"
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
	"github.com/alexfalkowski/go-service/internal/test"
	"github.com/alexfalkowski/go-service/module"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/telemetry"
	st "github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/token"
	"github.com/alexfalkowski/go-service/transport"
	geh "github.com/alexfalkowski/go-service/transport/events/http"
	"github.com/alexfalkowski/go-service/transport/http"
	"github.com/open-feature/go-sdk/openfeature"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	h "github.com/standard-webhooks/standard-webhooks/libraries/go"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func TestRunWithServer(t *testing.T) {
	Convey("Given I have valid configuration", t, func() {
		So(os.SetVariable("IN_CONFIG_FILE", test.Path("configs/config.yml")), ShouldBeNil)

		Convey("When I try to run an application that will shutdown in a second", func() {
			command := cmd.New("1.0.0")
			command.AddServer("server", "Start the server.", opts()...)
			command.RegisterInput(command.Root(), "env:IN_CONFIG_FILE")
			command.RegisterOutput(command.Root(), "env:OUT_CONFIG_FILE")
			command.SetArgs([]string{"server"})

			Convey("Then I should not see an error", func() {
				So(command.Run(), ShouldBeNil)
			})

			So(os.UnsetVariable("IN_CONFIG_FILE"), ShouldBeNil)
		})
	})
}

func TestRun(t *testing.T) {
	Convey("Given I have valid configuration", t, func() {
		So(os.SetVariable("CONFIG_FILE", test.Path("configs/config.yml")), ShouldBeNil)

		Convey("When I try to run an application that will shutdown in a second", func() {
			command := cmd.New("1.0.0")
			command.AddServer("server", "Start the server.", opts()...)
			command.RegisterInput(command.Root(), "env:CONFIG_FILE")

			Convey("Then I should not see an error", func() {
				So(command.Run(), ShouldBeNil)
			})

			So(os.UnsetVariable("CONFIG_FILE"), ShouldBeNil)
		})
	})
}

func TestInvalid(t *testing.T) {
	configs := []string{
		test.FilePath("configs/invalid_http.config.yml"),
		test.FilePath("configs/invalid_grpc.config.yml"),
		test.FilePath("configs/invalid_debug.config.yml"),
	}

	for _, config := range configs {
		Convey("When I try to run an application", t, func() {
			command := cmd.New("1.0.0")
			command.AddServer("server", "Start the server.", opts()...)
			command.RegisterInput(command.Root(), "env:CONFIG_FILE")
			command.SetArgs([]string{"server", "--input", config})

			Convey("Then I should not see an error", func() {
				err := command.Run()

				So(err, ShouldBeError)
				So(err.Error(), ShouldContainSubstring, "unknown port")
			})
		})
	}
}

func TestDisabled(t *testing.T) {
	Convey("When I try to run an application", t, func() {
		command := cmd.New("1.0.0")
		command.AddServer("server", "Start the server.", opts()...)
		command.RegisterInput(command.Root(), "env:CONFIG_FILE")
		command.SetArgs([]string{"server", "-i", test.FilePath("configs/disabled.config.yml")})

		Convey("Then I should see an error", func() {
			So(command.Run(), ShouldBeNil)
		})
	})
}

func TestExitOnRun(t *testing.T) {
	Convey("Given I have invalid configuration", t, func() {
		So(os.SetVariable("CONFIG_FILE", test.Path("configs/invalid_http.config.yml")), ShouldBeNil)

		Convey("When I try to run an application", func() {
			command := cmd.New("1.0.0")
			command.AddServer("server", "Start the server.", opts()...)
			command.RegisterInput(command.Root(), "env:CONFIG_FILE")
			command.SetArgs([]string{"server"})

			var exitCode int

			os.Exit = func(code int) {
				exitCode = code
			}

			command.ExitOnError()

			Convey("Then it should exit with a code of 1", func() {
				So(exitCode, ShouldEqual, 1)
			})

			So(os.UnsetVariable("CONFIG_FILE"), ShouldBeNil)
		})
	})
}

func TestClient(t *testing.T) {
	Convey("When I try to run a client", t, func() {
		opts := []fx.Option{fx.NopLogger}

		command := cmd.New("1.0.0")
		command.AddClient("client", "Start the client.", opts...)
		command.RegisterInput(command.Root(), "env:CONFIG_FILE")
		command.SetArgs([]string{"client"})

		Convey("Then I should not see an error", func() {
			So(command.Run(), ShouldBeNil)
		})
	})
}

func TestInvalidClient(t *testing.T) {
	configs := []string{
		test.Path("configs/invalid_http.config.yml"),
		test.Path("configs/invalid_grpc.config.yml"),
	}

	for _, config := range configs {
		Convey("Given I have invalid configuration", t, func() {
			So(os.SetVariable("TEST_CONFIG_FILE", config), ShouldBeNil)

			Convey("When I try to run an application", func() {
				command := cmd.New("1.0.0")
				command.AddClient("client", "Start the client.", opts()...)
				command.RegisterInput(command.Root(), "env:CONFIG_FILE")
				command.SetArgs([]string{"client", "--input", "env:TEST_CONFIG_FILE"})

				Convey("Then I should see an error", func() {
					err := command.Run()

					So(err, ShouldBeError)
					So(err.Error(), ShouldContainSubstring, "unknown port")
				})

				So(os.UnsetVariable("TEST_CONFIG_FILE"), ShouldBeNil)
			})
		})
	}
}

func registrations(logger *zap.Logger, cfg *http.Config, ua env.UserAgent, tracer trace.Tracer, _ env.Version) health.Registrations {
	if cfg == nil {
		return nil
	}

	t := 5 * time.Second
	nc := checker.NewNoopChecker()
	nr := server.NewRegistration("noop", t, nc)
	rt, _ := http.NewRoundTripper(http.WithClientLogger(logger), http.WithClientTracer(tracer), http.WithClientUserAgent(ua))
	hc := checker.NewHTTPChecker("https://google.com", rt, t)
	hr := server.NewRegistration("http", t, hc)

	return health.Registrations{nr, hr}
}

func healthObserver(healthServer *server.Server) (*shh.HealthObserver, error) {
	return &shh.HealthObserver{Observer: healthServer.Observe("noop")}, nil
}

func livenessObserver(healthServer *server.Server) *shh.LivenessObserver {
	return &shh.LivenessObserver{Observer: healthServer.Observe("noop")}
}

func readinessObserver(healthServer *server.Server) *shh.ReadinessObserver {
	return &shh.ReadinessObserver{Observer: healthServer.Observe("http")}
}

func grpcObserver(healthServer *server.Server) *shg.Observer {
	return &shg.Observer{Observer: healthServer.Observe("http")}
}

func invokeCache(_ cc.Cache) {}

func configs(_ *pg.Config, _ *feature.Config, _ *id.Config) {}

func meter(_ metric.Meter) {}

func featureClient(_ *openfeature.Client) {}

func webHooks(_ *h.Webhook, _ *geh.Receiver) {}

func environment(_ env.Name, _ env.UserAgent, _ env.Version) {}

func netTime(n st.Network) {
	_, _ = n.Now()
}

func crypt(signer *argon2.Signer, _ *ed25519.Signer, _ *rsa.Cipher, _ *aes.Cipher, _ *hmac.Signer, _ *ssh.Signer) error {
	msg := "hello"

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

func tokens(_ token.KID, _ *token.JWT, _ *token.Paseto, _ *token.Token) {}

func shutdown(s fx.Shutdowner) {
	go func(s fx.Shutdowner) {
		time.Sleep(time.Second)

		_ = s.Shutdown()
	}(s)
}

func opts() []fx.Option {
	return []fx.Option{
		module.Module, cmd.Module, config.Module, sd.Module,
		feature.Module, transport.Module, telemetry.Module, health.Module,
		sql.Module, hooks.Module, cache.Module, token.Module,
		fx.Provide(registrations), fx.Provide(healthObserver), fx.Provide(livenessObserver),
		fx.Provide(readinessObserver), fx.Provide(grpcObserver),
		fx.Invoke(shutdown), fx.Invoke(featureClient), fx.Invoke(webHooks), fx.Invoke(configs),
		fx.Invoke(meter), fx.Invoke(netTime), fx.Invoke(invokeCache),
		fx.Invoke(crypt), fx.Invoke(environment), fx.Invoke(tokens),
	}
}
