package http_test

import (
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/internal/test"
	nh "github.com/alexfalkowski/go-service/net/http"
	. "github.com/smartystreets/goconvey/convey"
)

func TestServerClose(t *testing.T) {
	Convey("When we check the error for a normal server close", t, func() {
		err := nh.ServerError(http.ErrServerClosed)

		Convey("Then there should be no error", func() {
			So(err, ShouldBeNil)
		})
	})

	Convey("When we check the error for abnormal server close", t, func() {
		err := nh.ServerError(test.ErrFailed)

		Convey("Then there should be an error", func() {
			So(err, ShouldBeError)
		})
	})
}
