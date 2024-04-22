package security_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/security"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestTLSMissingConfig(t *testing.T) {
	Convey("When I try to create with missing config", t, func() {
		c, err := security.NewTLSConfig(nil)

		Convey("Then I should have a default TLS config", func() {
			So(c, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
	})
}

func TestTLSInvalidPathConfig(t *testing.T) {
	Convey("When I try to create with missing config", t, func() {
		c, err := security.NewTLSConfig(&security.Config{Enabled: true})

		Convey("Then I should have a default TLS config", func() {
			So(c, ShouldBeNil)
			So(err, ShouldBeError)
		})
	})
}
