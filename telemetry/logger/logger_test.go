package logger_test

import (
	"context"
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestLogger(t *testing.T) {
	Convey("When I have a logger", t, func() {
		lc := fxtest.NewLifecycle(t)
		log := test.NewLogger(lc, test.NewTextLoggerConfig())

		Convey("Then I should log info", func() {
			log.Log(t.Context(), logger.NewText("test"), logger.Bool("yes", true))
		})

		Convey("Then I should log error", func() {
			log.Log(t.Context(), logger.NewMessage("test", context.Canceled), logger.Bool("yes", true))
		})
	})
}

func TestInvalidLogger(t *testing.T) {
	Convey("Given I have an invalid configuration", t, func() {
		lc := fxtest.NewLifecycle(t)
		cfg := &logger.Config{Kind: "wrong", Level: "debug"}
		params := logger.LoggerParams{
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
