package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/config"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestWatcher(t *testing.T) {
	Convey("Given I have valid configuration", t, func() {
		os.Setenv("CONFIG_FILE", "../test/config.yml")

		Convey("When I watch the config file", func() {
			lc := fxtest.NewLifecycle(t)
			sh := test.NewShutdowner()
			logger := test.NewLogger(lc)

			err := config.Watch(config.WatchParams{Lifecycle: lc, Shutdowner: sh, Logger: logger})
			So(err, ShouldBeNil)

			lc.RequireStart()

			Convey("Then I should shutdown when the config changes", func() {
				bytes, err := os.ReadFile(os.Getenv("CONFIG_FILE"))
				So(err, ShouldBeNil)

				config.WriteFileToEnv("CONFIG_FILE", bytes)
				So(err, ShouldBeNil)

				// Wait till we shutdown.
				time.Sleep(config.RandomWaitTime + time.Second)

				So(sh.Called(), ShouldBeTrue)
			})

			lc.RequireStop()
		})

		So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
	})
}
