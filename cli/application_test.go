package cli_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/cli"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx"
)

func TestApplicationRunWithServer(t *testing.T) {
	Convey("Given I have valid configuration", t, func() {
		t.Setenv("IN_CONFIG_FILE", test.Path("configs/config.yml"))

		Convey("When I try to run an application that will shutdown in a second", func() {
			app := cli.NewApplication(
				func(c cli.Commander) {
					cmd := c.AddServer("server", "Start the server.", test.Options()...)
					cmd.AddInput("env:IN_CONFIG_FILE")
					cmd.AddOutput("env:OUT_CONFIG_FILE")
				},
				cli.WithApplicationName(test.Name),
				cli.WithApplicationVersion(test.Version),
			)

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
			var exitCode int

			exit := func(code int) {
				exitCode = code
			}

			app := cli.NewApplication(
				func(c cli.Commander) {
					cmd := c.AddServer("server", "Start the server.", test.Options()...)
					cmd.AddInput("env:CONFIG_FILE")
				},
				cli.WithApplicationExit(exit),
			)

			app.ExitOnError(t.Context(), test.Name.String(), "server")

			Convey("Then it should exit with a code of 1", func() {
				So(exitCode, ShouldEqual, 1)
			})
		})
	})
}

func TestApplicationRun(t *testing.T) {
	Convey("Given I have valid configuration", t, func() {
		t.Setenv("CONFIG_FILE", test.Path("configs/config.yml"))

		Convey("When I try to run an application that will shutdown in a second", func() {
			app := cli.NewApplication(
				func(c cli.Commander) {
					cmd := c.AddServer("server", "Start the server.", test.Options()...)
					cmd.AddInput("env:CONFIG_FILE")
				},
				cli.WithApplicationName(test.Name),
				cli.WithApplicationVersion(test.Version),
				cli.WithApplicationExit(test.Exit),
			)

			Convey("Then I should see an error", func() {
				So(app.Run(t.Context()), ShouldBeError)
			})
		})
	})
}

func TestApplicationRunWithInvalidFlag(t *testing.T) {
	Convey("Given I have valid configuration", t, func() {
		t.Setenv("IN_CONFIG_FILE", test.Path("configs/config.yml"))

		Convey("When I try to run the application with an invalid flag", func() {
			app := cli.NewApplication(
				func(c cli.Commander) {
					cmd := c.AddServer("server", "Start the server.", test.Options()...)
					cmd.AddInput("env:IN_CONFIG_FILE")
					cmd.AddOutput("env:OUT_CONFIG_FILE")
				},
				cli.WithApplicationName(test.Name),
				cli.WithApplicationVersion(test.Version),
				cli.WithApplicationExit(test.Exit),
			)

			Convey("Then I should see an error", func() {
				So(app.Run(t.Context(), test.Name.String(), "server", "--invalid-flag"), ShouldBeError)
			})
		})
	})

	Convey("Given I have valid configuration", t, func() {
		t.Setenv("IN_CONFIG_FILE", test.Path("configs/config.yml"))

		Convey("When I try to run the application with an invalid flag", func() {
			app := cli.NewApplication(
				func(c cli.Commander) {
					cmd := c.AddClient("client", "Start the client.", test.Options()...)
					cmd.AddInput("env:IN_CONFIG_FILE")
					cmd.AddOutput("env:OUT_CONFIG_FILE")
				},
				cli.WithApplicationName(test.Name),
				cli.WithApplicationVersion(test.Version),
				cli.WithApplicationExit(test.Exit),
			)

			Convey("Then I should see an error", func() {
				So(app.Run(t.Context(), test.Name.String(), "client", "--invalid-flag"), ShouldBeError)
			})
		})
	})
}

func TestApplicationRunWithInvalidParams(t *testing.T) {
	Convey("Given I have valid configuration", t, func() {
		t.Setenv("IN_CONFIG_FILE", test.Path("configs/config.yml"))

		Convey("When I try to run an application that will shutdown in a second", func() {
			app := cli.NewApplication(
				func(c cli.Commander) {
					cmd := c.AddServer("server", "Start the server.", test.Options()...)
					cmd.AddInput("env:IN_CONFIG_FILE")
					cmd.AddOutput("env:OUT_CONFIG_FILE")
				},
				cli.WithApplicationName(test.Name),
				cli.WithApplicationVersion(test.Version),
				cli.WithApplicationExit(test.Exit),
			)

			Convey("Then I should not see an error", func() {
				So(app.Run(t.Context(), test.Name.String(), "server"), ShouldBeNil)
			})
		})
	})
}

func TestApplicationInvalid(t *testing.T) {
	configs := []string{
		test.FilePath("configs/invalid_http.config.yml"),
		test.FilePath("configs/invalid_grpc.config.yml"),
		test.FilePath("configs/invalid_debug.config.yml"),
	}

	for _, config := range configs {
		Convey("When I try to run an application", t, func() {
			app := cli.NewApplication(
				func(c cli.Commander) {
					cmd := c.AddServer("server", "Start the server.", test.Options()...)
					cmd.AddInput("env:CONFIG_FILE")
				},
				cli.WithApplicationName(test.Name),
				cli.WithApplicationVersion(test.Version),
				cli.WithApplicationExit(test.Exit),
			)

			Convey("Then I should not see an error", func() {
				err := app.Run(t.Context(), test.Name.String(), "server", "--input", config)

				So(err, ShouldBeError)
				So(err.Error(), ShouldContainSubstring, "unknown port")
			})
		})
	}
}

func TestApplicationDisabled(t *testing.T) {
	Convey("When I try to run an application", t, func() {
		app := cli.NewApplication(
			func(c cli.Commander) {
				cmd := c.AddServer("server", "Start the server.", test.Options()...)
				cmd.AddInput("env:CONFIG_FILE")
			},
			cli.WithApplicationName(test.Name),
			cli.WithApplicationVersion(test.Version),
			cli.WithApplicationExit(test.Exit),
		)

		Convey("Then I should see an error", func() {
			So(app.Run(t.Context(), test.Name.String(), "server", "-i", test.FilePath("configs/disabled.config.yml")), ShouldBeNil)
		})
	})
}

func TestApplicationClient(t *testing.T) {
	Convey("When I try to run a client", t, func() {
		opts := []fx.Option{fx.NopLogger}
		app := cli.NewApplication(
			func(c cli.Commander) {
				cmd := c.AddClient("client", "Start the client.", opts...)
				cmd.AddInput("env:CONFIG_FILE")
			},
			cli.WithApplicationName(test.Name),
			cli.WithApplicationVersion(test.Version),
			cli.WithApplicationExit(test.Exit),
		)

		Convey("Then I should not see an error", func() {
			So(app.Run(t.Context(), test.Name.String(), "client"), ShouldBeNil)
		})
	})
}

func TestApplicationInvalidClient(t *testing.T) {
	configs := []string{
		test.Path("configs/invalid_http.config.yml"),
		test.Path("configs/invalid_grpc.config.yml"),
	}

	for _, config := range configs {
		Convey("Given I have invalid configuration", t, func() {
			t.Setenv("TEST_CONFIG_FILE", config)

			Convey("When I try to run an application", func() {
				app := cli.NewApplication(
					func(c cli.Commander) {
						cmd := c.AddClient("client", "Start the client.", test.Options()...)
						cmd.AddInput("env:CONFIG_FILE")
					},
					cli.WithApplicationName(test.Name),
					cli.WithApplicationVersion(test.Version),
					cli.WithApplicationExit(test.Exit),
				)

				Convey("Then I should see an error", func() {
					err := app.Run(t.Context(), test.Name.String(), "client", "--input", "env:TEST_CONFIG_FILE")

					So(err, ShouldBeError)
					So(err.Error(), ShouldContainSubstring, "unknown port")
				})
			})
		})
	}
}
