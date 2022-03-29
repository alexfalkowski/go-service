package jaeger_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/logger/zap"
	"github.com/alexfalkowski/go-service/trace/opentracing/jaeger"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestJaeger(t *testing.T) {
	Convey("Given I have a configuration", t, func() {
		cfg := &jaeger.Config{
			Host: "localhost:6831",
		}

		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		Convey("When I register the trace system", func() {
			lc := fxtest.NewLifecycle(t)
			err := jaeger.Register(lc, logger, cfg)

			lc.RequireStart()

			Convey("Then I should have registered successfully", func() {
				So(err, ShouldBeNil)
			})

			lc.RequireStop()
		})
	})
}
