package http_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestGet(t *testing.T) {
	Convey("Given I have all the servers", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 100)))
		world.Start()

		world.ServeMux.HandleFunc("GET /hello", func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte("hello!"))
		})

		Convey("When I query for a greet", func() {
			_, _, err := world.Request(context.Background(), "http", http.MethodGet, "hello", http.Header{}, http.NoBody)
			So(err, ShouldBeNil)

			res, body, err := world.Request(context.Background(), "http", http.MethodGet, "hello", http.Header{}, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid greet", func() {
				So(res.StatusCode, ShouldEqual, http.StatusOK)
				So(body, ShouldEqual, "hello!")
			})

			world.Stop()
		})
	})
}

func TestLimiter(t *testing.T) {
	for _, f := range []string{"user-agent", "ip"} {
		Convey("Given I have a all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLimiter(test.NewLimiterConfig(f, "1s", 0)))
			world.Start()

			world.ServeMux.HandleFunc("GET /hello", func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte("hello!"))
			})

			Convey("When I query for a greet", func() {
				_, _, err := world.Request(context.Background(), "http", http.MethodGet, "hello", http.Header{}, http.NoBody)
				So(err, ShouldBeNil)

				res, _, err := world.Request(context.Background(), "http", http.MethodGet, "hello", http.Header{}, http.NoBody)
				So(err, ShouldBeNil)

				Convey("Then I should have been rate limited", func() {
					So(res.StatusCode, ShouldEqual, http.StatusTooManyRequests)
					So(res.Header.Get("Ratelimit"), ShouldNotBeBlank)
				})

				world.Stop()
			})
		})
	}
}

func TestClosedLimiter(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 100)))
		world.Start()

		ctx := context.Background()

		err := world.Server.Limiter.Close(ctx)
		So(err, ShouldBeNil)

		world.ServeMux.HandleFunc("GET /hello", func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte("hello!"))
		})

		Convey("When I query for a greet", func() {
			res, _, err := world.Request(context.Background(), "http", http.MethodGet, "hello", http.Header{}, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have an internal error", func() {
				So(res.StatusCode, ShouldEqual, http.StatusInternalServerError)
			})

			world.Stop()
		})
	})
}
