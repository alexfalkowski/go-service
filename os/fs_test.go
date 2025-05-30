package os_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	. "github.com/smartystreets/goconvey/convey"
)

func TestReadFile(t *testing.T) {
	for _, path := range []string{"none"} {
		Convey("When I check the path", t, func() {
			_, err := test.FS.ReadFile(path)

			Convey("Then it should not exist", func() {
				So(test.FS.IsNotExist(err), ShouldBeTrue)
				So(test.FS.PathExists(path), ShouldBeFalse)
			})
		})
	}
}

func TestPathExtension(t *testing.T) {
	for _, f := range []string{"file.yaml", "file.test.yaml", "test/.config/existing.client.yaml"} {
		Convey("When I try to get the extension of the file", t, func() {
			e := test.FS.PathExtension(f)

			Convey("Then the extension should be yaml", func() {
				So(e, ShouldEqual, "yaml")
			})
		})
	}

	Convey("When I try to get the extension of the file", t, func() {
		e := test.FS.PathExtension("file")

		Convey("Then the extension should be yaml", func() {
			So(e, ShouldBeEmpty)
		})
	})
}

func TestReadSource(t *testing.T) {
	t.Setenv("DUMMY", "yes")

	values := []*test.KeyValue[string, string]{
		{Key: "env:DUMMY", Value: "yes"},
		{Key: test.FilePath("configs/invalid.yml"), Value: "not:\n  our:\n    config: test"},
		{Key: "none", Value: "none"},
	}

	for _, value := range values {
		Convey("When I check the source", t, func() {
			data, err := test.FS.ReadSource(value.Key)

			Convey("Then I should have data", func() {
				So(err, ShouldBeNil)
				So(bytes.String(data), ShouldEqual, value.Value)
			})
		})
	}
}
