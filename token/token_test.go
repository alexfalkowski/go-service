package token_test

import (
	"context"
	"errors"
	"testing"

	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/token"
	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestGenerate(t *testing.T) {
	for _, kind := range []string{"key", "token"} {
		Convey("Given I have a invalid key token", t, func() {
			token := token.NewToken(test.NewToken(kind, "secrets/none"), test.Name, nil, nil)

			Convey("When I try to generate", func() {
				_, _, err := token.Generate(context.Background())

				Convey("Then I should have an error", func() {
					So(err, ShouldBeError)
				})
			})
		})
	}

	Convey("Given I have an invalid configuration", t, func() {
		token := token.NewToken(test.NewToken("none", "secrets/key"), test.Name, nil, nil)

		Convey("When I try to generate", func() {
			_, token, err := token.Generate(context.Background())

			Convey("Then I should have no token", func() {
				So(token, ShouldBeNil)
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given I have an missing configuration", t, func() {
		token := token.NewToken(nil, test.Name, nil, nil)

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
	for _, kind := range []string{"key", "token"} {
		Convey("Given I have a invalid key token", t, func() {
			token := token.NewToken(test.NewToken(kind, "secrets/none"), test.Name, nil, nil)

			Convey("When I try to verify", func() {
				_, err := token.Verify(context.Background(), nil)

				Convey("Then I should have an error", func() {
					So(err, ShouldBeError)
				})
			})
		})

		Convey("Given I have a valid key token", t, func() {
			token := token.NewToken(test.NewToken(kind, "secrets/"+kind), test.Name, nil, nil)

			Convey("When I try to verify", func() {
				_, err := token.Verify(context.Background(), []byte{})

				Convey("Then I should have an error", func() {
					So(err, ShouldBeError)
				})
			})
		})
	}

	Convey("Given I have an invalid configurionat", t, func() {
		token := token.NewToken(test.NewToken("none", "secrets/key"), test.Name, nil, nil)

		Convey("When I try to verify", func() {
			_, err := token.Verify(context.Background(), []byte{})

			Convey("Then I should have no error", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given I have an missing configurionat", t, func() {
		token := token.NewToken(nil, test.Name, nil, nil)

		Convey("When I try to verify", func() {
			_, err := token.Verify(context.Background(), []byte{})

			Convey("Then I should have no error", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestToken(t *testing.T) {
	Convey("When I generate a token", t, func() {
		key := token.Generate(test.Name)

		Convey("Then I should have a token", func() {
			So(key, ShouldNotBeBlank)
			So(token.Verify(test.Name, key), ShouldBeNil)
		})
	})

	keys := []string{
		"",
		"none_test_test",
		string(test.Name) + "_" + uuid.NewString() + "_test",
		string(test.Name) + "_" + uuid.NewString() + "_1",
		string(test.Name) + "_test_test",
		string(test.Name) + "_test_1",
	}

	for _, key := range keys {
		Convey("When I generate a token", t, func() {
			Convey("Then I should have an error", func() {
				err := token.Verify(test.Name, key)

				So(err, ShouldBeError)
				So(errors.Is(err, token.ErrInvalidMatch), ShouldBeTrue)
			})
		})
	}
}
