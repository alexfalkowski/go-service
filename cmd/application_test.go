package cmd_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/internal/test"
	"github.com/alexfalkowski/go-service/os"
	. "github.com/smartystreets/goconvey/convey"
)

func TestApplicationRunWithServer(t *testing.T) {
	Convey("Given I have valid configuration", t, func() {
		t.Setenv("IN_CONFIG_FILE", test.Path("configs/config.yml"))

		Convey("When I try to run an application that will shutdown in a second", func() {
			app := cmd.NewApplication(func(command *cmd.Command) {
				flags := command.AddServer("server", "Start the server.", opts()...)
				flags.AddInput("env:IN_CONFIG_FILE")
				flags.AddOutput("env:OUT_CONFIG_FILE")
			})

			Convey("Then I should not see an error", func() {
				So(app.Run(t.Context(), test.Name.String(), "server"), ShouldBeNil)
			})
		})
	})
}

func TestApplicationExitOnRun(t *testing.T) {
	Convey("Given I have invalid configuration", t, func() {
		t.Setenv("CONFIG_FILE", test.Path("configs/invalid_http.config.yml"))

		Convey("When I try to run an application", func() {
			app := cmd.NewApplication(func(command *cmd.Command) {
				flags := command.AddServer("server", "Start the server.", opts()...)
				flags.AddInput("env:CONFIG_FILE")
			})

			var exitCode int

			os.Exit = func(code int) {
				exitCode = code
			}

			app.ExitOnError(t.Context(), test.Name.String(), "server")

			Convey("Then it should exit with a code of 1", func() {
				So(exitCode, ShouldEqual, 1)
			})
		})
	})
}
