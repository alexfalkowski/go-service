package logger_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/internal/test"
	"github.com/alexfalkowski/go-service/transport/grpc/telemetry/logger"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc/codes"
)

func TestLogger(t *testing.T) {
	Convey("Given I have a logger", t, func() {
		lc := fxtest.NewLifecycle(t)
		l := test.NewLogger(lc)

		lc.RequireStart()

		Convey("When I try to get a logger func with a code", func() {
			f := logger.CodeToLogFunc(codes.DeadlineExceeded, l)

			Convey("Then I should a valid logger func", func() {
				So(f, ShouldNotBeNil)
			})
		})

		lc.RequireStop()
	})
}
