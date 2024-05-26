package limiter_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/transport/meta"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestValidLimiter(t *testing.T) {
	Convey("Given I have a valid format", t, func() {
		limiter.RegisterKey("user-agent", meta.UserAgent)

		c := &limiter.Config{Kind: "user-agent", Pattern: "0-S"}

		Convey("When I try to create a limiter", func() {
			_, _, err := limiter.New(c)

			Convey("Then I should have a valid limiter", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestInvalidLimiter(t *testing.T) {
	Convey("Given I have an invalid format", t, func() {
		limiter.RegisterKey("user-agent", meta.UserAgent)

		c := &limiter.Config{Kind: "user-agent", Pattern: "bob"}

		Convey("When I try to create a limiter", func() {
			_, _, err := limiter.New(c)

			Convey("Then I should have an invalid limiter", func() {
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I have an missing key", t, func() {
		c := &limiter.Config{Kind: "user-agent", Pattern: "bob"}

		Convey("When I try to create a limiter", func() {
			_, _, err := limiter.New(c)

			Convey("Then I should have an invalid limiter", func() {
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I have a disabled config", t, func() {
		limiter.RegisterKey("user-agent", meta.UserAgent)

		Convey("When I try to create a limiter", func() {
			c, _, err := limiter.New(nil)

			Convey("Then I should have an invalid limiter", func() {
				So(err, ShouldBeNil)
				So(c, ShouldBeNil)
			})
		})
	})
}
