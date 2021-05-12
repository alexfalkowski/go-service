package datadog_test

import (
	"os"
	"testing"

	"github.com/alexfalkowski/go-service/pkg/trace/opentracing/datadog"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestDatadog(t *testing.T) {
	Convey("Given I have a configuration", t, func() {
		os.Setenv("APP_NAME", "test")

		cfg, err := datadog.NewConfig()
		So(err, ShouldBeNil)

		Convey("When I register the trace system", func() {
			lc := fxtest.NewLifecycle(t)
			err := datadog.Register(lc, cfg)

			lc.RequireStart()

			Convey("Then I should have registered successfully", func() {
				So(err, ShouldBeNil)
			})

			lc.RequireStop()
		})

		So(os.Unsetenv("APP_NAME"), ShouldBeNil)
	})
}
