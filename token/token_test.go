package token_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/token/jwt"
	"github.com/alexfalkowski/go-service/v2/token/paseto"
	"github.com/alexfalkowski/go-service/v2/token/ssh"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGenerate(t *testing.T) {
	for _, kind := range []string{"jwt", "paseto", "none"} {
		cfg := test.NewToken(kind)
		ec := test.NewEd25519()
		signer, _ := ed25519.NewSigner(test.PEM, ec)
		verifier, _ := ed25519.NewVerifier(test.PEM, ec)
		gen := &id.UUID{}
		params := token.TokenParams{
			Config: cfg,
			Name:   test.Name,
			JWT: jwt.NewToken(jwt.TokenParams{
				Config:    cfg.JWT,
				Signer:    signer,
				Verifier:  verifier,
				Generator: gen,
			}),
			Paseto: paseto.NewToken(paseto.TokenParams{
				Config:    cfg.Paseto,
				Signer:    signer,
				Verifier:  verifier,
				Generator: gen,
			}),
		}
		tkn := token.NewToken(params)

		Convey("When I try to generate", t, func() {
			_, err := tkn.Generate("hello", test.UserID.String())

			Convey("Then I should have no error", func() {
				So(err, ShouldBeNil)
			})
		})
	}
}

//nolint:funlen
func TestVerify(t *testing.T) {
	for _, kind := range []string{"jwt", "paseto"} {
		cfg := test.NewToken(kind)
		ec := test.NewEd25519()
		signer, _ := ed25519.NewSigner(test.PEM, ec)
		verifier, _ := ed25519.NewVerifier(test.PEM, ec)
		gen := &id.UUID{}
		params := token.TokenParams{
			Config: cfg,
			Name:   test.Name,
			JWT: jwt.NewToken(jwt.TokenParams{
				Config:    cfg.JWT,
				Signer:    signer,
				Verifier:  verifier,
				Generator: gen,
			}),
			Paseto: paseto.NewToken(paseto.TokenParams{
				Config:    cfg.Paseto,
				Signer:    signer,
				Verifier:  verifier,
				Generator: gen,
			}),
		}
		tkn := token.NewToken(params)

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
		params := token.TokenParams{
			Config: cfg,
			Name:   test.Name,
			SSH:    ssh.NewToken(test.FS, cfg.SSH),
		}
		tkn := token.NewToken(params)

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

func TestWithNoConfig(t *testing.T) {
	Convey("When I try to create a token with no config", t, func() {
		params := token.TokenParams{Name: test.Name}
		token := token.NewToken(params)

		Convey("Then I should have no token", func() {
			So(token, ShouldBeNil)
		})
	})
}
