package os_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/os"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestPathExtension(t *testing.T) {
	for _, f := range []string{"file.yaml", "file.test.yaml", "test/.config/existing.client.yaml"} {
		Convey("When I try to get the extension of the file", t, func() {
			e := os.PathExtension(f)

			Convey("Then the extension should be yaml", func() {
				So(e, ShouldEqual, "yaml")
			})
		})
	}
}
