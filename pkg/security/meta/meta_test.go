package meta_test

import (
	"context"
	"testing"

	"github.com/alexfalkowski/go-service/pkg/security/meta"
	"github.com/form3tech-oss/jwt-go"
	. "github.com/smartystreets/goconvey/convey"
)

func TestWithToken(t *testing.T) {
	Convey("Given I have a valid token", t, func() {
		token := &jwt.Token{}
		ctx := meta.WithToken(context.Background(), token)

		Convey("When I try to get the token", func() {
			t := meta.Token(ctx)

			Convey("Then I should have a valid token", func() {
				So(t, ShouldEqual, token)
			})
		})
	})
}

func TestWithoutToken(t *testing.T) {
	Convey("Given I don't have a valid token", t, func() {
		Convey("When I try to get the token", func() {
			token := meta.Token(context.Background())

			Convey("Then I should have a missing token", func() {
				So(token, ShouldBeNil)
			})
		})
	})
}
