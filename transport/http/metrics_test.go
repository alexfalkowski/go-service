package http_test

import (
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/crypto/ed25519"
	"github.com/alexfalkowski/go-service/id"
	"github.com/alexfalkowski/go-service/internal/test"
	"github.com/alexfalkowski/go-service/token"
	"github.com/alexfalkowski/go-service/token/jwt"
	"github.com/alexfalkowski/go-service/types/ptr"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPrometheusAuthHTTP(t *testing.T) {
	Convey("Given I register the metrics handler", t, func() {
		cfg := test.NewToken("jwt")
		ec := test.NewEd25519()
		signer, _ := ed25519.NewSigner(ec)
		verifier, _ := ed25519.NewVerifier(ec)
		params := token.Params{
			Config: cfg,
			Name:   test.Name,
			JWT:    jwt.NewToken(cfg.JWT, signer, verifier, &id.UUID{}),
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

		err = world.Get(t.Context(), "not_existent", ptr)
		So(err, ShouldBeError)

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
