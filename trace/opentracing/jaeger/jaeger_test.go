package jaeger_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/trace/opentracing/jaeger"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestJaeger(t *testing.T) {
	Convey("Given I have a configuration", t, func() {
		cfg := test.NewJaegerConfig()
		lc := fxtest.NewLifecycle(t)

		Convey("When I register the trace system", func() {
			params := jaeger.TracerParams{Lifecycle: lc, Name: "test", Host: cfg.Host}
			_, err := jaeger.NewTracer(params)

			lc.RequireStart()

			Convey("Then I should have registered successfully", func() {
				So(err, ShouldBeNil)
			})

			lc.RequireStop()
		})
	})
}

func TestInvalidJaeger(t *testing.T) {
	Convey("Given I have a configuration", t, func() {
		lc := fxtest.NewLifecycle(t)

		Convey("When I register the trace system", func() {
			params := jaeger.TracerParams{Lifecycle: lc, Name: "test", Host: "invalid_host"}
			_, err := jaeger.NewTracer(params)

			lc.RequireStart()

			Convey("Then I should have registered unsuccessfully", func() {
				So(err, ShouldNotBeNil)
			})

			lc.RequireStop()
		})
	})
}
