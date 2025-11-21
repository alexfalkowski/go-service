package attributes_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAttributes(t *testing.T) {
	Convey("Then I should have attributes", t, func() {
		So(attributes.DBSystem("test").Value.AsString(), ShouldEqual, "test")
		So(attributes.Bool("test", true).Value.AsBool(), ShouldBeTrue)
		So(attributes.Float64("test", 0).Value.AsFloat64(), ShouldEqual, 0)
		So(attributes.Int64("test", 0).Value.AsInt64(), ShouldEqual, 0)
		So(attributes.String("test", "test").Value.AsString(), ShouldEqual, "test")
	})
}
