package feature_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/feature"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/open-feature/go-sdk/openfeature"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestFlags(t *testing.T) {
	Convey("Given I do not have a feature provider", t, func() {
		lc := fxtest.NewLifecycle(t)

		feature.Register(feature.ProviderParams{
			Lifecycle:      lc,
			Name:           test.Name,
			MetricProvider: test.NewPrometheusMeterProvider(lc),
		})

		client := feature.NewClient(test.Name)

		lc.RequireStart()

		Convey("When I get a flag", func() {
			attrs := map[string]any{"favorite_color": "blue"}
			v, err := client.BooleanValue(t.Context(), "v2_enabled", false, openfeature.NewEvaluationContext("tim@apple.com", attrs))
			So(err, ShouldBeNil)

			Convey("Then I should have missing flag", func() {
				So(v, ShouldBeFalse)
			})
		})

		lc.RequireStop()
	})

	Convey("Given I have a feature provider", t, func() {
		lc := fxtest.NewLifecycle(t)

		feature.Register(feature.ProviderParams{
			Lifecycle:       lc,
			Name:            test.Name,
			MetricProvider:  test.NewPrometheusMeterProvider(lc),
			FeatureProvider: openfeature.NoopProvider{},
		})

		client := feature.NewClient(test.Name)

		lc.RequireStart()

		Convey("When I get a flag", func() {
			attrs := map[string]any{"favorite_color": "blue"}
			v, err := client.BooleanValue(t.Context(), "v2_enabled", false, openfeature.NewEvaluationContext("tim@apple.com", attrs))
			So(err, ShouldBeNil)

			Convey("Then I should have missing flag", func() {
				So(v, ShouldBeFalse)
			})
		})

		lc.RequireStop()
	})
}
