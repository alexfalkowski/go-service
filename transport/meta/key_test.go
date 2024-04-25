package meta_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/transport/meta"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestKey(t *testing.T) {
	Convey("When I try to get an non existent key", t, func() {
		k := meta.NewKey(&limiter.Config{Enabled: true, Kind: "bob"})

		Convey("Then I should not get a key back", func() {
			So(k, ShouldBeNil)
		})
	})
}
