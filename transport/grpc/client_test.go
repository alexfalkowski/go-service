package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/security"
	"github.com/alexfalkowski/go-service/transport/grpc"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestClient(t *testing.T) {
	Convey("Given I have invalid creds", t, func() {
		c := &security.Config{Enabled: true, Cert: "bob", Key: "bob"}

		Convey("When I create an option", func() {
			_, err := grpc.WithClientSecure(c)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I have secure creds", t, func() {
		c := &security.Config{Enabled: true}

		Convey("When I create an option", func() {
			_, err := grpc.WithClientSecure(c)

			Convey("Then I should not have an error", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
