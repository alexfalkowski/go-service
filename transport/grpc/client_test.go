package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/crypto/tls"
	"github.com/alexfalkowski/go-service/transport/grpc"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestClient(t *testing.T) {
	Convey("Given I have invalid credentials", t, func() {
		c := &tls.Config{Cert: "bob", Key: "bob"}

		Convey("When I create a client", func() {
			_, err := grpc.NewClient("none", grpc.WithClientTLS(c))

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I have secure credentials", t, func() {
		c := &tls.Config{}

		Convey("When I create a client", func() {
			_, err := grpc.NewClient("none", grpc.WithClientTLS(c))

			Convey("Then I should not have an error", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
