package feature_test

import (
	"context"
	"testing"

	"github.com/alexfalkowski/go-service/client"
	"github.com/alexfalkowski/go-service/feature"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/transport/grpc/telemetry/tracer"
	"github.com/open-feature/go-sdk/openfeature"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

func TestFlipt(t *testing.T) {
	Convey("Given I have a flipt client", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		t, err := tracer.NewTracer(tracer.Params{Lifecycle: lc, Config: test.NewOTLPTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		m := test.NewMeter(lc)
		cfg := &feature.Config{Kind: "flipt", Config: client.Config{Host: "localhost:9000", Retry: test.NewRetry()}}
		p := feature.ClientParams{Config: cfg, Logger: logger, Tracer: t, Meter: m}

		c, err := feature.NewClient(p)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I get a missing flag", func() {
			attrs := map[string]any{"favorite_color": "blue"}
			_, err := c.BooleanValue(context.Background(), "v2_enabled", false, openfeature.NewEvaluationContext("tim@apple.com", attrs))

			Convey("Then I should have missing flag", func() {
				So(err, ShouldBeError)
			})
		})

		lc.RequireStop()
	})
}

func TestNoop(t *testing.T) {
	Convey("Given I have a flipt client", t, func() {
		p := feature.ClientParams{Config: &feature.Config{}}

		c, err := feature.NewClient(p)
		So(err, ShouldBeNil)

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
