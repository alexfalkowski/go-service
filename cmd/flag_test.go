package cmd_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/cmd"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

type tuple [2]string

func TestSplitFlag(t *testing.T) {
	tuples := []tuple{{"file", "file.yaml"}, {"file", "file.test.yaml"}, {"file", "test/.config/existing.client.yaml"}}

	for _, tu := range tuples {
		Convey("Given I have a flag", t, func() {
			Convey("When I split the flag", func() {
				k, l := cmd.SplitFlag(tu[0] + ":" + tu[1])

				Convey("Then I should have a valid split", func() {
					So(k, ShouldEqual, tu[0])
					So(l, ShouldEqual, tu[1])
				})
			})
		})
	}
}

func TestReadWriter(t *testing.T) {
	tuples := []tuple{{"file", "file.yaml"}, {"file", "file.test.yaml"}, {"file", "test/.config/existing.client.yaml"}}

	for _, tu := range tuples {
		Convey("Given I have a flag", t, func() {
			Convey("When I get a read writer", func() {
				rw := cmd.NewReadWriter(tu[0], tu[1])

				Convey("Then tI should have a valid split", func() {
					So(rw.Kind(), ShouldEqual, "yaml")
				})
			})
		})
	}
}
