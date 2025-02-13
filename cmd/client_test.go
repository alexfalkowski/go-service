package cmd_test

import (
	"context"
	"testing"

	"github.com/alexfalkowski/go-service/cmd"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

func TestStart(t *testing.T) {
	Convey("When I start a client command without an error", t, func() {
		lc := fxtest.NewLifecycle(t)
		cmd.Start(lc, func(_ context.Context) {})

		Convey("Then I should not have an error", func() {
			lc.RequireStart()
		})
	})

	Convey("When I start a client command with an error", t, func() {
		lc := fxtest.NewLifecycle(t)
		cmd.Start(lc, func(_ context.Context) { panic("whoops") })

		err := lc.Start(t.Context())

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})
}
