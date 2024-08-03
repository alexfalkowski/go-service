package feature_test

import (
	"context"
	"testing"

	"github.com/alexfalkowski/go-service/feature"
	"github.com/alexfalkowski/go-service/test"
	"github.com/open-feature/go-sdk/openfeature"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

func TestNoop(t *testing.T) {
	Convey("Given I have a flipt client", t, func() {
		lc := fxtest.NewLifecycle(t)

		feature.Register(feature.ProviderParams{
			Lifecycle:      lc,
			Name:           test.Name,
			MetricProvider: test.NewPrometheusMeterProvider(lc),
		})

		c := feature.NewClient(test.Name)

		lc.RequireStart()

		Convey("When I get a flag", func() {
			attrs := map[string]any{"favorite_color": "blue"}
			v, err := c.BooleanValue(context.Background(), "v2_enabled", false, openfeature.NewEvaluationContext("tim@apple.com", attrs))
			So(err, ShouldBeNil)

			Convey("Then I should have missing flag", func() {
				So(v, ShouldBeFalse)
			})
		})

		lc.RequireStop()
	})
}
