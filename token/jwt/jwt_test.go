package jwt_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/token/jwt"
	. "github.com/smartystreets/goconvey/convey"
)

func TestValid(t *testing.T) {
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(test.PEM, ec)
	verifier, _ := ed25519.NewVerifier(test.PEM, ec)

	Convey("When I generate a JWT token", t, func() {
		cfg := test.NewToken("jwt")
		token := jwt.NewToken(cfg.JWT, signer, verifier, &id.UUID{})

		tkn, err := token.Generate("aud")
		So(err, ShouldBeNil)

		Convey("Then I should have a token", func() {
			So(tkn, ShouldNotBeBlank)
		})

		Convey("Then I should be able to verify the token", func() {
			sub, err := token.Verify(tkn, "aud")
			So(err, ShouldBeNil)

			So(sub, ShouldEqual, "sub")
		})
	})
}

//nolint:funlen
func TestInvalid(t *testing.T) {
	cfg := test.NewToken("jwt")
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(test.PEM, ec)
	verifier, _ := ed25519.NewVerifier(test.PEM, ec)
	token := jwt.NewToken(cfg.JWT, signer, verifier, &id.UUID{})

	tokens := []string{
		"invalid",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
	}

	for _, tkn := range tokens {
		Convey("When I verify an invalid token", t, func() {
			_, err := token.Verify(tkn, "aud")

			Convey("Then I should have a error", func() {
				So(err, ShouldBeError)
			})
		})
	}

	Convey("When I generate a JWT token with invalid aud", t, func() {
		cfg := test.NewToken("jwt")
		token := jwt.NewToken(cfg.JWT, signer, verifier, &id.UUID{})

		tkn, err := token.Generate("aud")
		So(err, ShouldBeNil)

		token = jwt.NewToken(cfg.JWT, signer, verifier, &id.UUID{})

		Convey("Then I should have a token", func() {
			So(tkn, ShouldNotBeBlank)
		})

		Convey("Then I should have an error", func() {
			_, err := token.Verify(tkn, "test")
			So(err, ShouldBeError)
		})
	})

	Convey("When I generate a JWT token with invalid iss", t, func() {
		jcf := &jwt.Config{
			Subject:    "sub",
			Issuer:     "test",
			Expiration: "1h",
			KeyID:      "1234567890",
		}
		token := jwt.NewToken(jcf, signer, verifier, &id.UUID{})

		tkn, err := token.Generate("aud")
		So(err, ShouldBeNil)

		cfg := test.NewToken("jwt")
		token = jwt.NewToken(cfg.JWT, signer, verifier, &id.UUID{})

		Convey("Then I should have a token", func() {
			So(tkn, ShouldNotBeBlank)
		})

		Convey("Then I should have an error", func() {
			_, err := token.Verify(tkn, "aud")
			So(err, ShouldBeError)
		})
	})
}
