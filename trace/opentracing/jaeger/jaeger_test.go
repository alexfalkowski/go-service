package jaeger_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/logger/zap"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/trace/opentracing/jaeger"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestJaeger(t *testing.T) {
	Convey("Given I have a configuration", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		Convey("When I register the trace system", func() {
			lc := fxtest.NewLifecycle(t)
			_, err := jaeger.NewTracer(lc, logger, test.NewJaegerConfig())

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

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		Convey("When I register the trace system", func() {
			lc := fxtest.NewLifecycle(t)
			_, err := jaeger.NewTracer(lc, logger, &jaeger.Config{Host: "invalid_host"})

			lc.RequireStart()

			Convey("Then I should have registered unsuccessfully", func() {
				So(err, ShouldNotBeNil)
			})

			lc.RequireStop()
		})
	})
}
