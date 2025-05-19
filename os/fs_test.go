package os_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/internal/test"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFS(t *testing.T) {
	fs := test.FS

	for _, path := range []string{"none"} {
		Convey("When I check the path", t, func() {
			_, err := fs.ReadFile(path)

			Convey("Then it should not exist", func() {
				So(fs.IsNotExist(err), ShouldBeTrue)
				So(fs.PathExists(path), ShouldBeFalse)
			})
		})
	}

	for _, f := range []string{"file.yaml", "file.test.yaml", "test/.config/existing.client.yaml"} {
		Convey("When I try to get the extension of the file", t, func() {
			e := fs.PathExtension(f)

			Convey("Then the extension should be yaml", func() {
				So(e, ShouldEqual, "yaml")
			})
		})
	}
}
