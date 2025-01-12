package redis_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/cache/redis"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestRingOptions(t *testing.T) {
	Convey("When I try to create options with missing url", t, func() {
		_, err := redis.NewRingOptions(&redis.Config{})

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})

	Convey("When I try to create options with invalid url", t, func() {
		_, err := redis.NewRingOptions(&redis.Config{
			URL: test.Path("secrets/hooks"),
		})

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})
}
