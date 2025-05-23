package flag_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/cli/flag"
	"github.com/alexfalkowski/go-service/v2/strings"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSplitFlag(t *testing.T) {
	tuples := [][2]string{
		{"file", "file.yaml"},
		{"file", "file.test.yaml"},
		{"file", "test/.config/existing.client.yaml"},
	}

	for _, tuple := range tuples {
		Convey("When I split the flag", t, func() {
			k, l := flag.SplitFlag(strings.Join(":", tuple[0], tuple[1]))

			Convey("Then I should have a valid split", func() {
				So(k, ShouldEqual, tuple[0])
				So(l, ShouldEqual, tuple[1])
			})
		})
	}
}
