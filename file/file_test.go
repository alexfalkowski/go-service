package file_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/file"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestExtension(t *testing.T) {
	for _, f := range []string{"file.yaml", "file.test.yaml", "test/.config/existing.client.yaml"} {
		Convey("Given I have a file", t, func() {
			Convey("When I try to get the extension", func() {
				e := file.Extension(f)

				Convey("Then the extension should be yaml", func() {
					So(e, ShouldEqual, "yaml")
				})
			})
		})
	}
}
