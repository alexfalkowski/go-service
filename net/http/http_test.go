package http_test

import (
	"context"
	"net/http"
	"testing"

	nh "github.com/alexfalkowski/go-service/net/http"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
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

func TestWriteResponse(t *testing.T) {
	Convey("When we write with an erroneous writer", t, func() {
		w := &test.ErrResponseWriter{}
		ctx := context.Background()

		nh.WriteResponse(ctx, w, []byte("test"))

		Convey("Then we should record the error", func() {
			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})
}
