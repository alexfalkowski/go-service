package token_test

import (
	"context"
	"errors"
	"testing"

	"github.com/alexfalkowski/go-service/crypto/rand"
	"github.com/alexfalkowski/go-service/internal/test"
	"github.com/alexfalkowski/go-service/token"
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
	Convey("When I have a rand generator", t, func() {
		gen := rand.NewGenerator(rand.NewReader())

		Convey("When I generate a token", func() {
			key, err := token.Generate(test.Name, gen)
			So(err, ShouldBeNil)

			Convey("Then I should have a token", func() {
				So(key, ShouldNotBeBlank)
				So(token.Verify(test.Name, key), ShouldBeNil)
			})
		})
	})

	Convey("When I have a erroneous rand generator", t, func() {
		gen := rand.NewGenerator(&test.ErrReaderCloser{})

		Convey("When I generate a token", func() {
			_, err := token.Generate(test.Name, gen)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})

	keys := []string{
		"",
		"none_test_test",
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
