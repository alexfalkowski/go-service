package bytes_test

import (
	"strings"
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	eb "github.com/alexfalkowski/go-service/v2/encoding/bytes"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEncoder(t *testing.T) {
	Convey("Given I have bytes encoder", t, func() {
		encoder := eb.NewEncoder()

		buffer := test.Pool.Get()
		defer test.Pool.Put(buffer)

		Convey("When I encode a buffer", func() {
			err := encoder.Encode(buffer, bytes.NewBufferString("yes!"))
			So(err, ShouldBeNil)

			s := strings.TrimSpace(buffer.String())

			Convey("Then I should have valid message", func() {
				So(s, ShouldEqual, "yes!")
			})
		})

		Convey("When I encode an invalid type", func() {
			var str string
			err := encoder.Encode(buffer, &str)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I have bytes encoder", t, func() {
		encoder := eb.NewEncoder()

		Convey("When I decode a buffer", func() {
			var msg bytes.Buffer

			err := encoder.Decode(bytes.NewBufferString("test"), &msg)
			So(err, ShouldBeNil)

			Convey("Then I should have valid map", func() {
				So(msg.String(), ShouldEqual, "test")
			})
		})

		Convey("When I decode with an invalid type", func() {
			var msg string
			err := encoder.Decode(bytes.NewBufferString("test"), &msg)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
