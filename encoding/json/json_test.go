package json_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/alexfalkowski/go-service/encoding/json"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestEncoder(t *testing.T) {
	Convey("Given I have JSON encoder", t, func() {
		encoder := json.NewEncoder()

		bytes := test.Pool.Get()
		defer test.Pool.Put(bytes)

		msg := map[string]string{"test": "test"}

		Convey("When I encode the JSON", func() {
			err := encoder.Encode(bytes, msg)
			So(err, ShouldBeNil)

			s := strings.TrimSpace(bytes.String())

			Convey("Then I should have valid JSON", func() {
				So(s, ShouldEqual, "{\"test\":\"test\"}")
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
