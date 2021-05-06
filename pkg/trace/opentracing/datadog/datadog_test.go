package datadog_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/pkg/config"
	"github.com/alexfalkowski/go-service/pkg/logger/zap"
	"github.com/alexfalkowski/go-service/pkg/trace/opentracing/datadog"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestDatadog(t *testing.T) {
	Convey("Given I have a configuration", t, func() {
		cfg := &config.Config{AppName: "test"}
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc)
		So(err, ShouldBeNil)

		Convey("When I register the trace system", func() {
			err := datadog.Register(lc, cfg, logger)

			lc.RequireStart().RequireStop()

			Convey("Then I should have registered successfully", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
