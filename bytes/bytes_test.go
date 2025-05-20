package bytes_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/strings"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCopy(t *testing.T) {
	Convey("When I copy bytes", t, func() {
		hello := strings.Bytes("hello")
		helloCopy := bytes.Copy(hello)

		Convey("When I encode the YAML", func() {
			So(helloCopy, ShouldEqual, hello)
		})
	})
}
