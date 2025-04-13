package io_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alexfalkowski/go-service/config/io"
	"github.com/alexfalkowski/go-service/internal/test"
	. "github.com/smartystreets/goconvey/convey"
)

func TestReadWriter(t *testing.T) {
	Convey("When I get an existing read writer", t, func() {
		rw := io.NewReadWriter(test.Name, "file", test.Path("configs/config.yml"), test.FS)

		Convey("Then I should have a yml kind", func() {
			So(rw.Kind(), ShouldEqual, "yml")
		})
	})

	Convey("Given I have a config file in my home", t, func() {
		home, err := os.UserHomeDir()
		So(err, ShouldBeNil)

		path := filepath.Join(home, ".config", test.Name.String())

		err = os.MkdirAll(path, 0o777)
		So(err, ShouldBeNil)

		data, err := os.ReadFile(test.Path("configs/config.yml"))
		So(err, ShouldBeNil)

		err = os.WriteFile(filepath.Join(path, test.Name.String()+".yml"), data, 0o600)
		So(err, ShouldBeNil)

		Convey("When I get an missing read writer", func() {
			rw := io.NewReadWriter(test.Name, "file", test.Path("configs/config.yaml"), test.FS)

			Convey("Then I should have a yml kind", func() {
				So(rw.Kind(), ShouldEqual, "yml")
			})
		})

		err = os.RemoveAll(path)
		So(err, ShouldBeNil)
	})
}
