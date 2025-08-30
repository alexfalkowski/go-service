package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUnlimited(t *testing.T) {
	Convey("Given I have all the servers", t, func() {
		cfg := test.NewLimiterConfig("user-agent", "1s", 100)
		world := test.NewWorld(t,
			test.WithWorldTelemetry("otlp"),
			test.WithWorldClientLimiter(cfg),
			test.WithWorldServerLimiter(cfg),
			test.WithWorldHTTP(),
		)
		world.Register()
		world.HandleHello()
		world.RequireStart()

		Convey("When I query for a greet", func() {
			url := world.PathServerURL("http", "hello")

			_, _, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)
			So(err, ShouldBeNil)

			res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid greet", func() {
				So(res.StatusCode, ShouldEqual, http.StatusOK)
				So(body, ShouldEqual, "hello!")
			})

			world.RequireStop()
		})
	})
}

func TestServerLimiter(t *testing.T) {
	for _, f := range []string{"user-agent", "ip"} {
		Convey("Given I have a all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig(f, "1s", 0)), test.WithWorldHTTP())
			world.Register()
			world.HandleHello()
			world.RequireStart()

			Convey("When I query for a greet", func() {
				url := world.PathServerURL("http", "hello")

				_, _, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)
				So(err, ShouldBeNil)

				res, _, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)
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

func TestClientLimiter(t *testing.T) {
	for _, f := range []string{"user-agent", "ip"} {
		Convey("Given I have a all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldClientLimiter(test.NewLimiterConfig(f, "1s", 0)), test.WithWorldHTTP())
			world.Register()
			world.HandleHello()
			world.RequireStart()

			Convey("When I query for a greet", func() {
				url := world.PathServerURL("http", "hello")

				_, _, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)
				So(err, ShouldBeNil)

				_, _, err = world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)

				Convey("Then I should have been rate limited", func() {
					So(err, ShouldBeError)
					So(status.Code(err), ShouldEqual, http.StatusTooManyRequests)
				})

				world.RequireStop()
			})
		})
	}
}

func TestServerClosedLimiter(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())
		world.Register()
		world.HandleHello()
		world.RequireStart()

		err := world.Server.HTTPLimiter.Close(t.Context())
		So(err, ShouldBeNil)

		Convey("When I query for a greet", func() {
			url := world.PathServerURL("http", "hello")

			res, _, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have an internal error", func() {
				So(res.StatusCode, ShouldEqual, http.StatusInternalServerError)
			})

			world.RequireStop()
		})
	})
}

func TestClientClosedLimiter(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldClientLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())
		world.Register()
		world.HandleHello()
		world.RequireStart()

		Convey("When I query for a greet", func() {
			url := world.PathServerURL("http", "hello")

			err := world.Client.HTTPLimiter.Close(t.Context())
			So(err, ShouldBeNil)

			_, _, err = world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)

			Convey("Then I should have an internal error", func() {
				So(err, ShouldBeError)
				So(status.Code(err), ShouldEqual, http.StatusInternalServerError)
			})

			world.RequireStop()
		})
	})
}
