package slices_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/types/slices"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestAppendNil(t *testing.T) {
	for _, elem := range []*int{nil} {
		Convey("Given I have an empty array", t, func() {
			arr := []*int{}

			Convey("When I try to append a value", func() {
				arr = slices.Append(arr, elem)

				Convey("Then I should not have any elements", func() {
					So(arr, ShouldBeEmpty)
				})
			})
		})
	}
}

func TestAppendNotNil(t *testing.T) {
	integer := 2

	for _, elem := range []*int{&integer} {
		Convey("Given I have an empty array", t, func() {
			arr := []*int{}

			Convey("When I try to append a value", func() {
				arr = slices.Append(arr, elem)

				Convey("Then I should not have any elements", func() {
					So(arr, ShouldNotBeEmpty)
				})
			})
		})
	}
}
