package feature_test

import (
	"context"
	"testing"

	"github.com/alexfalkowski/go-service/feature"
	"github.com/open-feature/go-sdk/openfeature"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestFlipt(t *testing.T) {
	Convey("Given I have a flipt client", t, func() {
		c := feature.NewClient(&feature.Config{Kind: "flipt", Host: "localhost:9000"})

		Convey("When I get a missing flag", func() {
			attrs := map[string]any{"favorite_color": "blue"}
			_, err := c.BooleanValue(context.Background(), "v2_enabled", false, openfeature.NewEvaluationContext("tim@apple.com", attrs))

			Convey("Then I should have missing flag", func() {
				So(err, ShouldBeError)
			})
		})
	})
}

func TestNoop(t *testing.T) {
	Convey("Given I have a flipt client", t, func() {
		c := feature.NewClient(&feature.Config{})

		Convey("When I get a flag", func() {
			attrs := map[string]any{"favorite_color": "blue"}
			v, err := c.BooleanValue(context.Background(), "v2_enabled", false, openfeature.NewEvaluationContext("tim@apple.com", attrs))
			So(err, ShouldBeNil)

			Convey("Then I should have missing flag", func() {
				So(v, ShouldBeFalse)
			})
		})
	})
}
