package slices_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/types/ptr"
	"github.com/alexfalkowski/go-service/v2/types/slices"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEmptyAppendZero(t *testing.T) {
	for _, elem := range []*int{nil} {
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

	for _, elem := range []int{0} {
		Convey("Given I have an empty array", t, func() {
			arr := []int{}

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

func TestElemFunc(t *testing.T) {
	Convey("Given I have elems", t, func() {
		elems := []*string{ptr.Value("test")}

		Convey("When I try to get an existing elem", func() {
			elem, ok := slices.ElemFunc(elems, func(t *string) bool { return *t == "test" })

			Convey("Then I should have an elem", func() {
				So(elem, ShouldNotBeNil)
				So(ok, ShouldBeTrue)
			})
		})

		Convey("When I try to get a nonexistent elem", func() {
			elem, ok := slices.ElemFunc(elems, func(t *string) bool { return *t == "bob" })

			Convey("Then I should not have an elem", func() {
				So(elem, ShouldBeNil)
				So(ok, ShouldBeFalse)
			})
		})
	})
}
