package jaeger_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/pkg/config"
	opentracing "github.com/alexfalkowski/go-service/pkg/trace/opentracing/jaeger"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestJaeger(t *testing.T) {
	Convey("Given I have a configuration", t, func() {
		cfg := &config.Config{
			AppName: "test",
		}

		Convey("When I register the trace system", func() {
			lc := fxtest.NewLifecycle(t)
			err := opentracing.Register(lc, cfg)

			lc.RequireStart().RequireStop()

			Convey("Then I should have registered successfully", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
