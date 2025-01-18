package hooks_test

import (
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/transport/http/hooks"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	standardwebhooks "github.com/standard-webhooks/standard-webhooks/libraries/go"
)

func TestServer(t *testing.T) {
	Convey("Given I have an invalid writer", t, func() {
		handler := hooks.NewHandler(&standardwebhooks.Webhook{})
		writer := &test.BadResponseWriter{}
		req := &http.Request{Body: &test.BadReaderCloser{}}

		Convey("When I process a request", func() {
			handler.ServeHTTP(writer, req, nil)

			Convey("Then I should have a bad request", func() {
				So(writer.Code, ShouldEqual, 400)
			})
		})
	})
}

func TestClient(t *testing.T) {
	Convey("Given I have an invalid request body", t, func() {
		roundTripper := hooks.NewRoundTripper(&standardwebhooks.Webhook{}, nil)
		req := &http.Request{Body: &test.BadReaderCloser{}}

		Convey("When I process a request", func() {
			_, err := roundTripper.RoundTrip(req)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
