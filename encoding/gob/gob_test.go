package gob_test

import (
	"bytes"
	"encoding/base64"
	"testing"

	"github.com/alexfalkowski/go-service/encoding/gob"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestEncoder(t *testing.T) {
	t.Parallel()

	Convey("Given I have gob encoder", t, func() {
		encoder := gob.NewEncoder()

		bytes := test.Pool.Get()
		defer test.Pool.Put(bytes)

		msg := map[string]string{"test": "test"}

		Convey("When I encode the proto", func() {
			err := encoder.Encode(bytes, msg)
			So(err, ShouldBeNil)

			s := base64.StdEncoding.EncodeToString(bytes.Bytes())

			Convey("Then I should have valid proto", func() {
				So(s, ShouldEqual, "DX8EAQL/gAABDAEMAAAO/4AAAQR0ZXN0BHRlc3Q=")
			})
		})
	})

	Convey("Given I have gob encoder", t, func() {
		encoder := gob.NewEncoder()

		Convey("When I decode the gob", func() {
			m, err := base64.StdEncoding.DecodeString("DX8EAQL/gAABDAEMAAAO/4AAAQR0ZXN0BHRlc3Q=")
			So(err, ShouldBeNil)

			var msg map[string]string

			err = encoder.Decode(bytes.NewReader(m), &msg)
			So(err, ShouldBeNil)

			Convey("Then I should have valid msg", func() {
				So(msg, ShouldEqual, map[string]string{"test": "test"})
			})
		})
	})
}
