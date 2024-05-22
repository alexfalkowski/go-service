package feature_test

import (
	"context"
	"testing"

	"github.com/alexfalkowski/go-service/client"
	"github.com/alexfalkowski/go-service/feature"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/test"
	"github.com/open-feature/go-sdk/openfeature"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

func TestFlipt(t *testing.T) {
	Convey("Given I have a flipt client", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		tc := test.NewOTLPTracerConfig()
		tracer, err := tracer.NewTracer(lc, test.Environment, test.Version, tc, logger)
		So(err, ShouldBeNil)

		m := test.NewOTLPMeter(lc)
		cfg := &feature.Config{
			Kind:   "flipt",
			Config: &client.Config{Host: "localhost:9000", Retry: test.NewRetry(), Timeout: "5s"},
		}
		p := feature.ClientParams{Config: cfg, Logger: logger, Tracer: tracer, Meter: m}
		pr := feature.NewFeatureProvider(p)
		c := feature.NewClient(lc, pr)

		lc.RequireStart()

		Convey("When I get a missing flag", func() {
			attrs := map[string]any{"favorite_color": "blue"}
			_, err := c.BooleanValue(context.Background(), "v2_enabled", false, openfeature.NewEvaluationContext("tim@apple.com", attrs))

			Convey("Then I should have missing flag", func() {
				So(err, ShouldBeError)
				So(feature.IsNotFoundError(err), ShouldBeTrue)
			})
		})

		lc.RequireStop()
	})

	Convey("Given I have a flipt client", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		tc := test.NewOTLPTracerConfig()
		tracer, err := tracer.NewTracer(lc, test.Environment, test.Version, tc, logger)
		So(err, ShouldBeNil)

		m := test.NewOTLPMeter(lc)
		cfg := &feature.Config{
			Kind:   "flipt",
			Config: &client.Config{Host: "localhost:9000", Retry: test.NewRetry(), Timeout: "5s"},
		}
		p := feature.ClientParams{Config: cfg, Logger: logger, Tracer: tracer, Meter: m}
		pr := feature.NewFeatureProvider(p)
		c := feature.NewClient(lc, pr)

		lc.RequireStart()

		Convey("When I ping", func() {
			err := feature.Ping(context.Background(), c)

			Convey("Then I should be up", func() {
				So(err, ShouldBeNil)
			})
		})

		lc.RequireStop()
	})
}

func TestNoop(t *testing.T) {
	Convey("Given I have a flipt client", t, func() {
		lc := fxtest.NewLifecycle(t)
		p := feature.ClientParams{Config: &feature.Config{}}
		pr := feature.NewFeatureProvider(p)
		c := feature.NewClient(lc, pr)

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
