package tls_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/crypto/tls"
	"github.com/alexfalkowski/go-service/internal/test"
	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig(t *testing.T) {
	configs := []*tls.Config{nil, {}}

	for _, c := range configs {
		Convey("When I try to create with missing config", t, func() {
			c, err := tls.NewConfig(test.FS, c)

			Convey("Then I should have a default TLS config", func() {
				So(c, ShouldNotBeNil)
				So(err, ShouldBeNil)
			})
		})
	}

	configs = []*tls.Config{
		test.NewTLSConfig("certs/client-cert.pem", "secrets/none"),
		test.NewTLSConfig("secrets/none", "certs/client-key.pem"),
		test.NewTLSConfig("secrets/hooks", "certs/client-key.pem"),
	}

	for _, c := range configs {
		Convey("When I try to create with missing config", t, func() {
			_, err := tls.NewConfig(test.FS, c)

			Convey("Then I should have an errror", func() {
				So(err, ShouldBeError)
			})
		})
	}
}
