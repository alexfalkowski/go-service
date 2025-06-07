package logger_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/transport/grpc/telemetry/logger"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc/codes"
)

func TestLogger(t *testing.T) {
	Convey("When I try to get a level for a timeout", t, func() {
		level := logger.CodeToLevel(codes.DeadlineExceeded)

		Convey("Then I should an level of error", func() {
			So(level, ShouldEqual, logger.LevelError)
		})
	})
}
