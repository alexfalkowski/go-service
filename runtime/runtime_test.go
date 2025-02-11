package runtime_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/internal/test"
	"github.com/alexfalkowski/go-service/runtime"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestPanic(t *testing.T) {
	Convey("When I have an error", t, func() {
		f := func() { runtime.Must(test.ErrFailed) }

		Convey("Then I should panic", func() {
			So(f, ShouldPanic)
		})
	})

	Convey("When I don't have an error", t, func() {
		f := func() { runtime.Must(nil) }

		Convey("Then I should panic", func() {
			So(f, ShouldNotPanic)
		})
	})
}

func TestRecover(t *testing.T) {
	type fun func() (err error)

	errPanic := func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = runtime.ConvertRecover(r)
			}
		}()

		panic(test.ErrFailed)
	}

	strPanic := func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = runtime.ConvertRecover(r)
			}
		}()

		panic("test")
	}

	intPanic := func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = runtime.ConvertRecover(r)
			}
		}()

		panic(1)
	}

	for _, f := range []fun{errPanic, strPanic, intPanic} {
		Convey("When I panic in a function", t, func() {
			err := f()

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	}
}
