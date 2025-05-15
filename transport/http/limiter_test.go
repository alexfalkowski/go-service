package http_test

import (
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/internal/test"
	"github.com/alexfalkowski/go-service/strings"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGet(t *testing.T) {
	Convey("Given I have all the servers", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		world.HandleFunc("GET /hello", func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write(strings.Bytes("hello!"))
		})

		Convey("When I query for a greet", func() {
			_, _, err := world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodGet, "hello", http.Header{}, http.NoBody)
			So(err, ShouldBeNil)

			res, body, err := world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodGet, "hello", http.Header{}, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid greet", func() {
				So(res.StatusCode, ShouldEqual, http.StatusOK)
				So(body, ShouldEqual, "hello!")
			})

			world.RequireStop()
		})
	})
}

func TestLimiter(t *testing.T) {
	for _, f := range []string{"user-agent", "ip"} {
		Convey("Given I have a all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLimiter(test.NewLimiterConfig(f, "1s", 0)), test.WithWorldHTTP())
			world.Register()
			world.RequireStart()

			world.HandleFunc("GET /hello", func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write(strings.Bytes("hello!"))
			})

			Convey("When I query for a greet", func() {
				_, _, err := world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodGet, "hello", http.Header{}, http.NoBody)
				So(err, ShouldBeNil)

				res, _, err := world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodGet, "hello", http.Header{}, http.NoBody)
				So(err, ShouldBeNil)

				Convey("Then I should have been rate limited", func() {
					So(res.StatusCode, ShouldEqual, http.StatusTooManyRequests)
					So(res.Header.Get("Ratelimit"), ShouldNotBeBlank)
				})

				world.RequireStop()
			})
		})
	}
}

func TestClosedLimiter(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		err := world.Limiter.Close(t.Context())
		So(err, ShouldBeNil)

		world.HandleFunc("GET /hello", func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write(strings.Bytes("hello!"))
		})

		Convey("When I query for a greet", func() {
			res, _, err := world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodGet, "hello", http.Header{}, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have an internal error", func() {
				So(res.StatusCode, ShouldEqual, http.StatusInternalServerError)
			})

			world.RequireStop()
		})
	})
}
