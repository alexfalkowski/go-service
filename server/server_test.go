package server_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/server"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

func TestServer(t *testing.T) {
	Convey("Given I have a server", t, func() {
		lc := fxtest.NewLifecycle(t)
		l := test.NewLogger(lc)
		sh := test.NewShutdowner()
		srv := &errServe{}
		s := server.NewServer("test", srv, l, sh)

		Convey("When I start", func() {
			s.Start()
			time.Sleep(1 * time.Second)

			Convey("Then it should shutdown", func() {
				So(sh.Called(), ShouldBeTrue)
			})
		})
	})
}

type errServe struct{}

func (e *errServe) IsEnabled() bool {
	return true
}

func (e *errServe) Serve() error {
	return os.ErrNotExist
}

func (e *errServe) Shutdown(_ context.Context) error {
	return nil
}

func (e *errServe) String() string {
	return "test"
}
