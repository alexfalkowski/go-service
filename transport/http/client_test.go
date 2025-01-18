package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/crypto/tls"
	"github.com/alexfalkowski/go-service/transport/http"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestClient(t *testing.T) {
	Convey("Given I have invalid creds", t, func() {
		c := &tls.Config{Cert: "bob", Key: "bob"}

		Convey("When I create a client", func() {
			_, err := http.NewClient(http.WithClientTLS(c))

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I have secure creds", t, func() {
		c := &tls.Config{}

		Convey("When I create a client", func() {
			_, err := http.NewClient(http.WithClientTLS(c))

			Convey("Then I should not have an error", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
