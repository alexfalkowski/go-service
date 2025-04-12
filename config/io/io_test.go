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
	tuples := [][2]string{
		{"file", "file.yaml"},
		{"file", "file.test.yaml"},
		{"file", "test/.config/existing.client.yaml"},
	}

	for _, tuple := range tuples {
		Convey("When I get a read writer", t, func() {
			rw := io.NewReadWriter(test.Name, tuple[0], tuple[1], test.FS)

			Convey("Then I should have a valid split", func() {
				So(rw.Kind(), ShouldEqual, "yaml")
			})
		})
	}

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

		Convey("When I try to read", func() {
			rw := io.NewCommon(test.Name, test.FS, &test.BadReaderWriter{})

			bytes, err := rw.Read()
			So(err, ShouldBeNil)

			Convey("Then I should have a valid config", func() {
				So(bytes, ShouldNotBeNil)
				So(rw.Kind(), ShouldEqual, "yml")
			})
		})

		err = os.RemoveAll(path)
		So(err, ShouldBeNil)
	})
}
