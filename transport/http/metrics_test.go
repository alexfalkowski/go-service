package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/types/ptr"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPrometheusHTTP(t *testing.T) {
	Convey("Given I register the metrics handler", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("prometheus"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())
		world.Register()

		_, err := world.OpenDatabase()
		So(err, ShouldBeNil)

		world.RequireStart()

		Convey("When I query metrics", func() {
			ctx, cancel := test.Timeout()
			defer cancel()

			header := http.Header{}
			url := world.NamedServerURL("http", "metrics")

			res, body, err := world.ResponseWithBody(ctx, url, http.MethodGet, header, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have valid metrics", func() {
				So(res.StatusCode, ShouldEqual, http.StatusOK)

				So(body, ShouldContainSubstring, "go_info")
				So(body, ShouldContainSubstring, `db_system="redis"`)
				So(body, ShouldContainSubstring, `db_system_name="pg"`)
				So(body, ShouldContainSubstring, "system")
				So(body, ShouldContainSubstring, "process")
				So(body, ShouldContainSubstring, "runtime")
			})
		})

		world.RequireStop()
	})
}

func TestPrometheusAuthHTTP(t *testing.T) {
	Convey("Given I register the metrics handler", t, func() {
		cfg := test.NewToken("jwt")
		ec := test.NewEd25519()
		signer, _ := ed25519.NewSigner(test.PEM, ec)
		verifier, _ := ed25519.NewVerifier(test.PEM, ec)
		gen := uuid.NewGenerator()
		tkn := token.NewToken(test.Name, cfg, test.FS, signer, verifier, gen)

		world := test.NewWorld(t,
			test.WithWorldTelemetry("prometheus"),
			test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)),
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
			url := world.NamedServerURL("http", "metrics")

			res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have valid metrics", func() {
				So(res.StatusCode, ShouldEqual, http.StatusOK)

				So(body, ShouldContainSubstring, "go_info")
				So(body, ShouldContainSubstring, `db_system="redis"`)
				So(body, ShouldContainSubstring, `db_system_name="pg"`)
				So(body, ShouldContainSubstring, "system")
				So(body, ShouldContainSubstring, "process")
				So(body, ShouldContainSubstring, "runtime")
			})
		})

		world.RequireStop()
	})
}
