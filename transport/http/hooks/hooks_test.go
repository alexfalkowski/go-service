package hooks_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/transport/http/hooks"
	. "github.com/smartystreets/goconvey/convey"
)

func TestVerify(t *testing.T) {
	Convey("Given I have an invalid request", t, func() {
		hook := hooks.NewWebhook(nil, nil)
		req := &http.Request{Body: &test.ErrReaderCloser{}}

		Convey("When I process a request", func() {
			err := hook.Verify(req)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}

func TestSign(t *testing.T) {
	Convey("Given I have an invalid request", t, func() {
		hook := hooks.NewWebhook(nil, nil)
		req := &http.Request{Body: &test.ErrReaderCloser{}}

		Convey("When I process a request", func() {
			err := hook.Sign(req)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}

func TestRoundTripper(t *testing.T) {
	Convey("Given I have an invalid request", t, func() {
		hook := hooks.NewWebhook(nil, nil)
		rt := hooks.NewRoundTripper(hook, nil)
		req := &http.Request{Body: &test.ErrReaderCloser{}}

		Convey("When I process a request", func() {
			_, err := rt.RoundTrip(req)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
