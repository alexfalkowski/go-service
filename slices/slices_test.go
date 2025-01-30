package slices_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/slices"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestEmptyAppendZero(t *testing.T) {
	var integer int

	for _, elem := range []*int{nil, &integer} {
		Convey("Given I have an empty array", t, func() {
			arr := []*int{}

			Convey("When I try to append a value", func() {
				arr = slices.AppendNotZero(arr, elem)

				Convey("Then I should not have any elements", func() {
					So(arr, ShouldBeEmpty)
				})
			})
		})
	}
}

func TestEmptyAppendNil(t *testing.T) {
	for _, elem := range []*int{nil} {
		Convey("Given I have an empty array", t, func() {
			arr := []*int{}

			Convey("When I try to append a value", func() {
				arr = slices.AppendNotNil(arr, elem)

				Convey("Then I should not have any elements", func() {
					So(arr, ShouldBeEmpty)
				})
			})
		})
	}
}

func TestAppendZero(t *testing.T) {
	integer := 2

	for _, elem := range []*int{&integer} {
		Convey("Given I have an empty array", t, func() {
			arr := []*int{}

			Convey("When I try to append a value", func() {
				arr = slices.AppendNotZero(arr, elem)

				Convey("Then I should not have any elements", func() {
					So(arr, ShouldNotBeEmpty)
				})
			})
		})
	}
}

func TestAppendNil(t *testing.T) {
	integer := 2

	for _, elem := range []*int{&integer} {
		Convey("Given I have an empty array", t, func() {
			arr := []*int{}

			Convey("When I try to append a value", func() {
				arr = slices.AppendNotNil(arr, elem)

				Convey("Then I should not have any elements", func() {
					So(arr, ShouldNotBeEmpty)
				})
			})
		})
	}
}
