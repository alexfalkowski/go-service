package http_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/crypto/ed25519"
	"github.com/alexfalkowski/go-service/crypto/rand"
	"github.com/alexfalkowski/go-service/id"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/token"
	"github.com/alexfalkowski/go-service/types/ptr"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestPrometheusAuthHTTP(t *testing.T) {
	Convey("Given I register the metrics handler", t, func() {
		kid, _ := token.NewKID(rand.NewGenerator(rand.NewReader()))
		a, _ := ed25519.NewSigner(test.NewEd25519())
		id := id.Default
		jwt := token.NewJWT(kid, a, id)
		pas := token.NewPaseto(a, id)
		token := token.NewToken(test.NewToken("jwt", "secrets/jwt"), test.Name, jwt, pas)

		world := test.NewWorld(t,
			test.WithWorldTelemetry("prometheus"),
			test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 100)),
			test.WithWorldToken(token, token),
		)
		world.Register()

		_, err := world.OpenDatabase()
		So(err, ShouldBeNil)

		world.RequireStart()

		ptr := ptr.Empty[string]()
		ctx := context.Background()

		err = world.Get(ctx, "not_existent", ptr)
		So(err, ShouldBeError)

		Convey("When I query metrics", func() {
			header := http.Header{}

			res, body, err := world.ResponseWithBody(ctx, "http", world.ServerHost(), http.MethodGet, "metrics", header, http.NoBody)
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

func TestPrometheusInsecureHTTP(t *testing.T) {
	Convey("Given I register the metrics handler", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("prometheus"), test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 100)))
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

func TestPrometheusSecureHTTP(t *testing.T) {
	Convey("Given I register the metrics handler", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("prometheus"), test.WithWorldSecure())
		world.Register()

		_, err := world.OpenDatabase()
		So(err, ShouldBeNil)

		world.RequireStart()

		Convey("When I query metrics", func() {
			header := http.Header{}

			res, body, err := world.ResponseWithBody(context.Background(), "https", world.ServerHost(), http.MethodGet, "metrics", header, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have valid metrics", func() {
				So(res.StatusCode, ShouldEqual, http.StatusOK)

				So(body, ShouldContainSubstring, "go_info")
				So(body, ShouldContainSubstring, "sql_max_open_total")
			})
		})

		world.RequireStop()
	})
}
