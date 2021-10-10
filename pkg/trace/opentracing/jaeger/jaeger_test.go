package jaeger_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/pkg/trace/opentracing/jaeger"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestJaeger(t *testing.T) {
	Convey("Given I have a configuration", t, func() {
		cfg := &jaeger.Config{
			Host: "localhost:6831",
		}

		Convey("When I register the trace system", func() {
			lc := fxtest.NewLifecycle(t)
			err := jaeger.Register(lc, cfg)

			lc.RequireStart()

			Convey("Then I should have registered successfully", func() {
				So(err, ShouldBeNil)
			})

			lc.RequireStop()
		})
	})
}
