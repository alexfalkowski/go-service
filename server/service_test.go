package server_test

import (
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/internal/test"
	"github.com/alexfalkowski/go-service/server"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestServer(t *testing.T) {
	Convey("Given I have an erroneous server", t, func() {
		lc := fxtest.NewLifecycle(t)
		l := test.NewLogger(lc, test.NewJSONLoggerConfig())
		sh := test.NewShutdowner()
		srv := &test.ErrServer{}
		svc := server.NewService("test", srv, l, sh)

		Convey("When I start", func() {
			svc.Start()
			time.Sleep(1 * time.Second)

			Convey("Then it should shutdown", func() {
				So(sh.Called(), ShouldBeTrue)
			})
		})
	})

	Convey("Given I have a server", t, func() {
		sh := test.NewShutdowner()
		srv := &test.NoopServer{}
		svc := server.NewService("test", srv, nil, sh)

		Convey("When I start", func() {
			svc.Start()
			time.Sleep(1 * time.Second)

			Convey("Then it should not shutdown", func() {
				So(sh.Called(), ShouldBeFalse)
			})

			svc.Stop(t.Context())
		})
	})
}
