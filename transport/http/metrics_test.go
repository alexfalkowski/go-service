package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/token/jwt"
	th "github.com/alexfalkowski/go-service/v2/transport/http"
	"github.com/alexfalkowski/go-service/v2/types/ptr"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	th.Register(test.FS)
}

func TestPrometheusAuthHTTP(t *testing.T) {
	Convey("Given I register the metrics handler", t, func() {
		cfg := test.NewToken("jwt")
		ec := test.NewEd25519()
		signer, _ := ed25519.NewSigner(test.PEM, ec)
		verifier, _ := ed25519.NewVerifier(test.PEM, ec)
		params := token.TokenParams{
			Config: cfg,
			Name:   test.Name,
			JWT: jwt.NewToken(jwt.TokenParams{
				Config:    cfg.JWT,
				Signer:    signer,
				Verifier:  verifier,
				Generator: &id.UUID{},
			}),
		}
		tkn := token.NewToken(params)

		world := test.NewWorld(t,
			test.WithWorldTelemetry("prometheus"),
			test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 100)),
			test.WithWorldToken(tkn, tkn),
			test.WithWorldHTTP(),
		)
		world.Register()

		_, err := world.OpenDatabase()
		So(err, ShouldBeNil)

		world.RequireStart()

		ptr := ptr.Zero[string]()

		ok, err := world.Get(t.Context(), "not_existent", ptr)
		So(ok, ShouldBeFalse)
		So(err, ShouldBeNil)

		Convey("When I query metrics", func() {
			header := http.Header{}

			res, body, err := world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodGet, "metrics", header, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have valid metrics", func() {
				So(res.StatusCode, ShouldEqual, http.StatusOK)

				So(body, ShouldContainSubstring, "go_info")
				So(body, ShouldContainSubstring, "cache_misses_total")
				So(body, ShouldContainSubstring, "sql_max_open_total")
				So(body, ShouldContainSubstring, "system")
				So(body, ShouldContainSubstring, "process")
				So(body, ShouldContainSubstring, "runtime")
			})
		})

		world.RequireStop()
	})
}

func TestPrometheusHTTP(t *testing.T) {
	Convey("Given I register the metrics handler", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("prometheus"), test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())
		world.Register()

		_, err := world.OpenDatabase()
		So(err, ShouldBeNil)

		world.RequireStart()

		Convey("When I query metrics", func() {
			ctx, cancel := test.Timeout()
			defer cancel()

			header := http.Header{}

			res, body, err := world.ResponseWithBody(ctx, "http", world.InsecureServerHost(), http.MethodGet, "metrics", header, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have valid metrics", func() {
				So(res.StatusCode, ShouldEqual, http.StatusOK)

				So(body, ShouldContainSubstring, "go_info")
				So(body, ShouldContainSubstring, "sql_max_open_total")
				So(body, ShouldContainSubstring, "system")
				So(body, ShouldContainSubstring, "process")
				So(body, ShouldContainSubstring, "runtime")
			})
		})

		world.RequireStop()
	})
}
