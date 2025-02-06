package os_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/os"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestFS(t *testing.T) {
	fs := os.NewFS()

	for _, path := range []string{"none"} {
		Convey("When I check the path", t, func() {
			_, err := fs.ReadFile(path)

			Convey("Then it should not exist", func() {
				So(fs.IsNotExist(err), ShouldBeTrue)
				So(fs.PathExists(path), ShouldBeFalse)
			})
		})
	}
}
