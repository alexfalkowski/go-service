package sync_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/sync"
	"github.com/alexfalkowski/go-service/types/ptr"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestArray(t *testing.T) {
	Convey("Given I have an empty array", t, func() {
		arr := sync.NewArray[string]()

		Convey("When I try to append a value", func() {
			arr.Add(ptr.Value("test"))

			Convey("Then I should elements", func() {
				So(arr.Elements(), ShouldNotBeEmpty)
			})
		})
	})
}
