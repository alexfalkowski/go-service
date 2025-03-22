package http_test

import (
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/crypto/ed25519"
	"github.com/alexfalkowski/go-service/id"
	"github.com/alexfalkowski/go-service/internal/test"
	"github.com/alexfalkowski/go-service/token"
	"github.com/alexfalkowski/go-service/types/ptr"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPrometheusAuthHTTP(t *testing.T) {
	Convey("Given I register the metrics handler", t, func() {
		cfg := test.NewToken("jwt", "secrets/jwt")
		kid := token.NewKID(cfg)
		signer, _ := ed25519.NewSigner(test.NewEd25519())
		params := token.Params{
			Config: cfg,
			Name:   test.Name,
			JWT:    token.NewJWT(kid, signer, &id.UUID{}),
		}
		token := token.NewToken(params)

		world := test.NewWorld(t,
			test.WithWorldTelemetry("prometheus"),
			test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 100)),
			test.WithWorldToken(token, token),
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

			res, body, err := world.ResponseWithBody(t.Context(), "http", world.ServerHost(), http.MethodGet, "metrics", header, http.NoBody)
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

			res, body, err := world.ResponseWithBody(ctx, "http", world.ServerHost(), http.MethodGet, "metrics", header, http.NoBody)
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
