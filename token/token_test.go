package token_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/token"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGenerate(t *testing.T) {
	for _, kind := range []string{"jwt", "paseto", "none"} {
		cfg := test.NewToken(kind)
		ec := test.NewEd25519()
		signer, _ := ed25519.NewSigner(test.PEM, ec)
		verifier, _ := ed25519.NewVerifier(test.PEM, ec)
		gen := &id.UUID{}
		tkn := token.NewToken(test.Name, cfg, test.FS, signer, verifier, gen)

		Convey("When I try to generate", t, func() {
			_, err := tkn.Generate("hello", test.UserID.String())

			Convey("Then I should have no error", func() {
				So(err, ShouldBeNil)
			})
		})
	}
}

func TestVerify(t *testing.T) {
	for _, kind := range []string{"jwt", "paseto"} {
		cfg := test.NewToken(kind)
		ec := test.NewEd25519()
		signer, _ := ed25519.NewSigner(test.PEM, ec)
		verifier, _ := ed25519.NewVerifier(test.PEM, ec)
		gen := &id.UUID{}
		tkn := token.NewToken(test.Name, cfg, test.FS, signer, verifier, gen)

		Convey("Given I generate a token", t, func() {
			gen, err := tkn.Generate("hello", test.UserID.String())
			So(err, ShouldBeNil)

			Convey("When I try to verify", func() {
				sub, err := tkn.Verify(gen, "hello")

				Convey("Then I should have no error", func() {
					So(err, ShouldBeNil)
					So(sub, ShouldEqual, test.UserID.String())
				})
			})
		})
	}

	for _, kind := range []string{"ssh", "none"} {
		cfg := test.NewToken(kind)
		tkn := token.NewToken(test.Name, cfg, test.FS, nil, nil, nil)

		Convey("Given I generate a token", t, func() {
			gen, err := tkn.Generate("", "")
			So(err, ShouldBeNil)

			Convey("When I try to verify", func() {
				_, err := tkn.Verify(gen, "")

				Convey("Then I should have no error", func() {
					So(err, ShouldBeNil)
				})
			})
		})
	}
}
