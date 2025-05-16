package os_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/internal/test"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPathExtension(t *testing.T) {
	for _, f := range []string{"file.yaml", "file.test.yaml", "test/.config/existing.client.yaml"} {
		Convey("When I try to get the extension of the file", t, func() {
			e := test.FS.PathExtension(f)

			Convey("Then the extension should be yaml", func() {
				So(e, ShouldEqual, "yaml")
			})
		})
	}
}
