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
	Convey("Given I have gob encoder", t, func() {
		e := gob.NewEncoder()

		b := test.Pool.Get()
		defer test.Pool.Put(b)

		msg := map[string]string{"test": "test"}

		Convey("When I encode the proto", func() {
			err := e.Encode(b, msg)
			So(err, ShouldBeNil)

			s := base64.StdEncoding.EncodeToString(b.Bytes())

			Convey("Then I should have valid proto", func() {
				So(s, ShouldEqual, "DX8EAQL/gAABDAEMAAAO/4AAAQR0ZXN0BHRlc3Q=")
			})
		})
	})

	Convey("Given I have gob encoder", t, func() {
		e := gob.NewEncoder()

		Convey("When I decode the gob", func() {
			m, err := base64.StdEncoding.DecodeString("DX8EAQL/gAABDAEMAAAO/4AAAQR0ZXN0BHRlc3Q=")
			So(err, ShouldBeNil)

			var msg map[string]string

			err = e.Decode(bytes.NewReader(m), &msg)
			So(err, ShouldBeNil)

			Convey("Then I should have valid msg", func() {
				So(msg, ShouldEqual, map[string]string{"test": "test"})
			})
		})
	})
}
