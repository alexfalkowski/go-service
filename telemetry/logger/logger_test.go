package logger_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestLogger(t *testing.T) {
	Convey("Given I have an invalid configuration", t, func() {
		lc := fxtest.NewLifecycle(t)
		cfg := &logger.Config{Kind: "wrong", Level: "debug"}
		params := logger.Params{
			Lifecycle:   lc,
			Config:      cfg,
			ID:          test.ID,
			Name:        test.Name,
			Version:     test.Version,
			Environment: test.Environment,
		}

		Convey("When I try to get a logger", func() {
			_, err := logger.NewLogger(params)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
