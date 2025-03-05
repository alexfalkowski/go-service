package limiter_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/transport/meta"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

//nolint:funlen
func TestLimiter(t *testing.T) {
	lc := fxtest.NewLifecycle(t)

	Convey("Given I have an missing key", t, func() {
		config := &limiter.Config{Kind: "user-agent", Tokens: 0, Interval: "1s"}

		Convey("When I try to create a limiter", func() {
			_, err := limiter.New(lc, config)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I have a disabled config", t, func() {
		limiter.RegisterKey("user-agent", meta.UserAgent)

		Convey("When I try to create a limiter", func() {
			limiter, err := limiter.New(lc, nil)

			Convey("Then I should have no limiter", func() {
				So(err, ShouldBeNil)
				So(limiter, ShouldBeNil)
			})
		})
	})

	Convey("Given I have a valid format", t, func() {
		limiter.RegisterKey("user-agent", meta.UserAgent)

		Convey("When I try to create a limiter", func() {
			config := &limiter.Config{Kind: "user-agent", Tokens: 0, Interval: "1s"}
			limiter, err := limiter.New(lc, config)

			Convey("Then I should have a limiter", func() {
				So(err, ShouldBeNil)
				So(limiter, ShouldNotBeNil)
			})

			err = limiter.Close(t.Context())
			So(err, ShouldBeNil)
		})
	})

	Convey("Given I have a limiter", t, func() {
		limiter.RegisterKey("user-agent", meta.UserAgent)

		config := &limiter.Config{Kind: "user-agent", Tokens: 0, Interval: "1s"}

		limiter, err := limiter.New(lc, config)
		So(err, ShouldBeNil)

		Convey("When I try take when the limiter is closed", func() {
			err = limiter.Close(t.Context())
			So(err, ShouldBeNil)

			_, _, err := limiter.Take(t.Context())

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
