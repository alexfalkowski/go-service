package cli_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/cli"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
	. "github.com/smartystreets/goconvey/convey"
)

func TestApplicationRun(t *testing.T) {
	Convey("Given I have valid configuration", t, func() {
		config := test.FilePath("configs/config.yml")

		os.Args = []string{test.Name.String(), "server", "-i", config}
		cli.Name = test.Name
		cli.Version = test.Version

		Convey("When I try to run an application that will shutdown in a second", func() {
			app := cli.NewApplication(
				func(c cli.Commander) {
					cmd := c.AddServer("server", "Start the server.", test.Options()...)
					cmd.AddInput(strings.Empty)
				},
			)

			Convey("Then I should not see an error", func() {
				So(app.Run(t.Context()), ShouldBeNil)
			})
		})
	})
}

func TestApplicationExitOnRun(t *testing.T) {
	Convey("Given I have invalid configuration", t, func() {
		config := test.FilePath("configs/invalid_http.config.yml")

		os.Args = []string{test.Name.String(), "server", "-i", config}

		Convey("When I try to run an application", func() {
			var exitCode int
			exit := func(code int) {
				exitCode = code
			}

			app := cli.NewApplication(
				func(c cli.Commander) {
					cmd := c.AddServer("server", "Start the server.", test.Options()...)
					cmd.AddInput(strings.Empty)
				},
				cli.WithApplicationExit(exit),
			)

			app.ExitOnError(t.Context())

			Convey("Then it should exit with a code of 1", func() {
				So(exitCode, ShouldEqual, 1)
			})
		})
	})
}

func TestApplicationRunWithInvalidFlag(t *testing.T) {
	Convey("When I try to run the server with an invalid flag", t, func() {
		os.Args = []string{test.Name.String(), "server", "--invalid-flag"}
		cli.Name = test.Name
		cli.Version = test.Version

		app := cli.NewApplication(
			func(c cli.Commander) {
				cmd := c.AddServer("server", "Start the server.", test.Options()...)
				cmd.AddInput(strings.Empty)
			},
		)

		Convey("Then I should see an error", func() {
			So(app.Run(t.Context()), ShouldBeError)
		})
	})

	Convey("When I try to run the client with an invalid flag", t, func() {
		os.Args = []string{test.Name.String(), "client", "--invalid-flag"}
		cli.Name = test.Name
		cli.Version = test.Version

		app := cli.NewApplication(
			func(c cli.Commander) {
				cmd := c.AddClient("client", "Start the client.", test.Options()...)
				cmd.AddInput(strings.Empty)
			},
		)

		Convey("Then I should see an error", func() {
			So(app.Run(t.Context()), ShouldBeError)
		})
	})
}

func TestApplicationRunWithInvalidParams(t *testing.T) {
	Convey("Given I have valid configuration", t, func() {
		config := test.FilePath("configs/config.yml")

		os.Args = []string{test.Name.String(), "server", "-i", config}
		cli.Name = test.Name
		cli.Version = test.Version

		Convey("When I try to run an application that will shutdown in a second", func() {
			app := cli.NewApplication(
				func(c cli.Commander) {
					cmd := c.AddServer("server", "Start the server.", test.Options()...)
					cmd.AddInput(strings.Empty)
				},
			)

			Convey("Then I should not see an error", func() {
				So(app.Run(t.Context()), ShouldBeNil)
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
			os.Args = []string{test.Name.String(), "server", "-i", config}
			cli.Name = test.Name
			cli.Version = test.Version

			app := cli.NewApplication(
				func(c cli.Commander) {
					cmd := c.AddServer("server", "Start the server.", test.Options()...)
					cmd.AddInput(strings.Empty)
				},
			)

			Convey("Then I should see an error", func() {
				err := app.Run(t.Context())

				So(err, ShouldBeError)
				So(err.Error(), ShouldContainSubstring, "unknown port")
			})
		})
	}
}

func TestApplicationDisabled(t *testing.T) {
	Convey("When I try to run an application", t, func() {
		os.Args = []string{test.Name.String(), "server", "-i", test.FilePath("configs/disabled.config.yml")}
		cli.Name = test.Name
		cli.Version = test.Version

		app := cli.NewApplication(
			func(c cli.Commander) {
				cmd := c.AddServer("server", "Start the server.", test.Options()...)
				cmd.AddInput(strings.Empty)
			},
		)

		Convey("Then I should not see an error", func() {
			So(app.Run(t.Context()), ShouldBeNil)
		})
	})
}

func TestApplicationClient(t *testing.T) {
	Convey("When I try to run a client", t, func() {
		os.Args = []string{test.Name.String(), "client"}
		cli.Name = test.Name
		cli.Version = test.Version

		opts := []di.Option{di.NoLogger}
		app := cli.NewApplication(
			func(c cli.Commander) {
				cmd := c.AddClient("client", "Start the client.", opts...)
				cmd.AddInput(strings.Empty)
			},
		)

		Convey("Then I should not see an error", func() {
			So(app.Run(t.Context()), ShouldBeNil)
		})
	})
}

func TestApplicationInvalidClient(t *testing.T) {
	configs := []string{
		test.FilePath("configs/invalid_http.config.yml"),
		test.FilePath("configs/invalid_grpc.config.yml"),
	}

	for _, config := range configs {
		Convey("When I try to run an application", t, func() {
			os.Args = []string{test.Name.String(), "client", "-i", config}
			cli.Name = test.Name
			cli.Version = test.Version

			app := cli.NewApplication(
				func(c cli.Commander) {
					cmd := c.AddClient("client", "Start the client.", test.Options()...)
					cmd.AddInput(strings.Empty)
				},
			)

			Convey("Then I should see an error", func() {
				err := app.Run(t.Context())

				So(err, ShouldBeError)
				So(err.Error(), ShouldContainSubstring, "unknown port")
			})
		})
	}
}
