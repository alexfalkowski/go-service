package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/transport/grpc"
	. "github.com/smartystreets/goconvey/convey"
)

func TestClient(t *testing.T) {
	grpc.Register(test.FS)

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
