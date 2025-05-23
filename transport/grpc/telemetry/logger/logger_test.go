package logger_test

import (
	"log/slog"
	"testing"

	"github.com/alexfalkowski/go-service/v2/transport/grpc/telemetry/logger"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc/codes"
)

func TestLogger(t *testing.T) {
	Convey("When I try to get a level", t, func() {
		level := logger.CodeToLevel(codes.DeadlineExceeded)

		Convey("Then I should a valid logger func", func() {
			So(level, ShouldEqual, slog.LevelError)
		})
	})
}
