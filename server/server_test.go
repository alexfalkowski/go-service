package server_test

import (
	"context"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/internal/test"
	"github.com/alexfalkowski/go-service/server"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

func TestServer(t *testing.T) {
	Convey("Given I have an erroneous server", t, func() {
		lc := fxtest.NewLifecycle(t)
		l := test.NewLogger(lc)
		sh := test.NewShutdowner()
		srv := &test.ErrServer{}
		server := server.NewServer("test", srv, l, sh)

		Convey("When I start", func() {
			server.Start()
			time.Sleep(1 * time.Second)

			Convey("Then it should shutdown", func() {
				So(sh.Called(), ShouldBeTrue)
			})
		})
	})

	Convey("Given I have a server", t, func() {
		sh := test.NewShutdowner()
		srv := &test.NoopServer{}
		server := server.NewServer("test", srv, nil, sh)

		Convey("When I start", func() {
			server.Start()
			time.Sleep(1 * time.Second)

			Convey("Then it should not shutdown", func() {
				So(sh.Called(), ShouldBeFalse)
			})

			server.Stop(context.Background())
		})
	})
}
