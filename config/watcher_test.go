package config_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/config"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestWatcher(t *testing.T) {
	runtimes := []string{"os", "container"}

	for _, runtime := range runtimes {
		Convey("Given I have valid configuration", t, func() {
			os.Setenv("CONFIG_FILE", "../test/config.yml")

			Convey(fmt.Sprintf("When I watch the config file for runtime %s", runtime), func() {
				lc := fxtest.NewLifecycle(t)
				sh := test.NewShutdowner()
				logger := test.NewLogger(lc)
				waitTime := config.WaitTime(time.Second)
				params := config.WatchParams{Lifecycle: lc, Shutdowner: sh, Logger: logger, WaitTime: waitTime, Config: &test.Config{Runtime: runtime}}

				err := config.Watch(params)
				So(err, ShouldBeNil)

				lc.RequireStart()

				Convey("Then I should shutdown when the config changes", func() {
					bytes, err := os.ReadFile(os.Getenv("CONFIG_FILE"))
					So(err, ShouldBeNil)

					config.WriteFileToEnv("CONFIG_FILE", bytes)
					So(err, ShouldBeNil)

					// Wait till we shutdown.
					time.Sleep(3 * time.Second)

					So(sh.Called(), ShouldBeTrue)
				})

				lc.RequireStop()
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	}
}
