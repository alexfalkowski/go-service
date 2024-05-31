package ssh_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/transport"
	ss "github.com/alexfalkowski/go-service/transport/ssh"
	"github.com/gliderlabs/ssh"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

type handler struct{}

func (*handler) Handle(_ ssh.Context, _ []string) error {
	return nil
}

func TestClient(t *testing.T) {
	Convey("Given I have an SSH server", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		p := ss.ServerParams{
			Logger:     logger,
			Config:     test.NewInsecureTransportConfig().SSH,
			Handler:    &handler{},
			Shutdowner: test.NewShutdowner(),
		}

		s, err := ss.NewServer(p)
		So(err, ShouldBeNil)

		transport.Register(transport.RegisterParams{Lifecycle: lc, Servers: []transport.Server{s}})

		lc.RequireStart()

		Convey("When I connect with a client", func() {
			c, err := ss.NewClient("localhost:"+p.Config.Port, ss.WithClientLogger(logger), ss.WithClientTimeout("10s"))
			So(err, ShouldBeNil)

			defer c.Close()

			r, err := c.Run("test")
			So(err, ShouldBeNil)

			Convey("Then I should have an error", func() {
				So(string(r), ShouldEqual, "test: successful")
			})
		})

		lc.RequireStop()
	})
}
