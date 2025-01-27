package token_test

import (
	"context"
	"testing"

	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/token"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestGenerate(t *testing.T) {
	Convey("When I generate a key", t, func() {
		key, err := token.GenerateKey()
		So(err, ShouldBeNil)

		Convey("Then I should have a key", func() {
			So(key, ShouldNotBeBlank)

			err := token.VerifyKey(key)
			So(err, ShouldBeNil)
		})
	})

	Convey("Given I have a invalid key token", t, func() {
		config := &token.Config{
			Kind:       "key",
			Key:        test.Path("secrets/none"),
			Subject:    "sub",
			Audience:   "aud",
			Issuer:     "iss",
			Expiration: "1h",
		}
		token := token.NewToken(config, nil, nil)

		Convey("When I try to generate", func() {
			_, _, err := token.Generate(context.Background())

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I have an invalid configuration", t, func() {
		token := token.NewToken(test.NewToken("none"), nil, nil)

		Convey("When I try to generate", func() {
			_, token, err := token.Generate(context.Background())

			Convey("Then I should have no token", func() {
				So(token, ShouldBeNil)
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given I have an missing configuration", t, func() {
		token := token.NewToken(nil, nil, nil)

		Convey("When I try to generate", func() {
			_, token, err := token.Generate(context.Background())

			Convey("Then I should have no token", func() {
				So(token, ShouldBeNil)
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestVerify(t *testing.T) {
	Convey("Given I have a invalid key token", t, func() {
		config := &token.Config{
			Kind:       "key",
			Key:        test.Path("secrets/none"),
			Subject:    "sub",
			Audience:   "aud",
			Issuer:     "iss",
			Expiration: "1h",
		}
		token := token.NewToken(config, nil, nil)

		Convey("When I try to verify", func() {
			_, err := token.Verify(context.Background(), nil)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I have a valid key token", t, func() {
		token := token.NewToken(test.NewToken("key"), nil, nil)

		Convey("When I try to verify", func() {
			_, err := token.Verify(context.Background(), []byte{})

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I have an invalid configurionat", t, func() {
		token := token.NewToken(test.NewToken("none"), nil, nil)

		Convey("When I try to verify", func() {
			_, err := token.Verify(context.Background(), []byte{})

			Convey("Then I should have no error", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given I have an missing configurionat", t, func() {
		token := token.NewToken(nil, nil, nil)

		Convey("When I try to verify", func() {
			_, err := token.Verify(context.Background(), []byte{})

			Convey("Then I should have no error", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
