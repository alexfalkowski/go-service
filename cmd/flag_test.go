package cmd_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

type tuple [2]string

func TestSplitFlag(t *testing.T) {
	tuples := []tuple{
		{"file", "file.yaml"},
		{"file", "file.test.yaml"},
		{"file", "test/.config/existing.client.yaml"},
	}

	for _, tuple := range tuples {
		Convey("When I split the flag", t, func() {
			k, l := cmd.SplitFlag(tuple[0] + ":" + tuple[1])

			Convey("Then I should have a valid split", func() {
				So(k, ShouldEqual, tuple[0])
				So(l, ShouldEqual, tuple[1])
			})
		})
	}
}

func TestReadWriter(t *testing.T) {
	tuples := []tuple{
		{"file", "file.yaml"},
		{"file", "file.test.yaml"},
		{"file", "test/.config/existing.client.yaml"},
	}

	for _, tuple := range tuples {
		Convey("When I get a read writer", t, func() {
			rw := cmd.NewReadWriter(tuple[0], tuple[1], test.FS)

			Convey("Then tI should have a valid split", func() {
				So(rw.Kind(), ShouldEqual, "yaml")
			})
		})
	}
}
