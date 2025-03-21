package logger_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/internal/test"
	"github.com/alexfalkowski/go-service/telemetry/logger"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestLogger(t *testing.T) {
	Convey("Given I have an invalid file system", t, func() {
		lc := fxtest.NewLifecycle(t)
		cfg := test.NewOTLPLoggerConfig()
		params := logger.Params{
			Lifecycle: lc, Config: cfg,
			Environment: test.Environment, Version: test.Version,
			FileSystem: &test.ErrFS{},
		}

		Convey("When I try to get a logger", func() {
			_, err := logger.NewLogger(params)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I have an invalid configuration", t, func() {
		lc := fxtest.NewLifecycle(t)
		cfg := &logger.Config{Kind: "wrong", Level: "debug"}
		params := logger.Params{
			Lifecycle: lc, Config: cfg,
			Environment: test.Environment, Version: test.Version,
			FileSystem: test.FS,
		}

		Convey("When I try to get a logger", func() {
			_, err := logger.NewLogger(params)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
