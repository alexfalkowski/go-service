package errors_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestServerClose(t *testing.T) {
	Convey("When we check the error for a normal server close", t, func() {
		err := errors.ServerError(http.ErrServerClosed)

		Convey("Then there should be no error", func() {
			So(err, ShouldBeNil)
		})
	})

	Convey("When we check the error for abnormal server close", t, func() {
		err := errors.ServerError(test.ErrFailed)

		Convey("Then there should be an error", func() {
			So(err, ShouldBeError)
		})
	})
}
