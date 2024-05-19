package tls_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/crypto/tls"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestConfig(t *testing.T) {
	cs := []*tls.Config{nil, {}}

	for _, c := range cs {
		Convey("When I try to create with missing config", t, func() {
			c, err := tls.NewConfig(c)

			Convey("Then I should have a default TLS config", func() {
				So(c, ShouldNotBeNil)
				So(err, ShouldBeNil)
			})
		})
	}
}
