package json_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/bytes"
	"github.com/alexfalkowski/go-service/encoding/json"
	"github.com/alexfalkowski/go-service/internal/test"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEncoder(t *testing.T) {
	Convey("Given I have JSON encoder", t, func() {
		encoder := json.NewEncoder()

		buffer := test.Pool.Get()
		defer test.Pool.Put(buffer)

		msg := map[string]string{"test": "test"}

		Convey("When I encode the JSON", func() {
			err := encoder.Encode(buffer, msg)
			So(err, ShouldBeNil)

			s := bytes.TrimSpace(test.Pool.Copy(buffer))

			Convey("Then I should have valid JSON", func() {
				So(bytes.String(s), ShouldEqual, "{\"test\":\"test\"}")
			})
		})
	})

	Convey("Given I have JSON encoder", t, func() {
		encoder := json.NewEncoder()

		Convey("When I decode the JSON", func() {
			var msg map[string]string

			err := encoder.Decode(bytes.NewReader([]byte("{\"test\":\"test\"}")), &msg)
			So(err, ShouldBeNil)

			Convey("Then I should have valid map", func() {
				So(msg, ShouldEqual, map[string]string{"test": "test"})
			})
		})
	})
}
