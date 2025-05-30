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
		params := token.Params{
			Config: cfg,
			Name:   test.Name,
			JWT:    jwt.NewToken(cfg.JWT, signer, verifier, gen),
			Paseto: paseto.NewToken(cfg.Paseto, signer, verifier, gen),
		}
		tkn := token.NewToken(params)

		Convey("When I try to generate", t, func() {
			_, _, err := tkn.Generate(t.Context())

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
		params := token.Params{
			Config: cfg,
			Name:   test.Name,
			JWT:    jwt.NewToken(cfg.JWT, signer, verifier, gen),
			Paseto: paseto.NewToken(cfg.Paseto, signer, verifier, gen),
		}
		tkn := token.NewToken(params)

		Convey("Given I generate a token", t, func() {
			_, gen, err := tkn.Generate(t.Context())
			So(err, ShouldBeNil)

			Convey("When I try to verify", func() {
				ctx, err := tkn.Verify(t.Context(), gen)

				Convey("Then I should have no error", func() {
					So(err, ShouldBeNil)
					So(token.Subject(ctx).String(), ShouldEqual, "sub")
				})
			})
		})
	}

	for _, kind := range []string{"ssh", "none"} {
		cfg := test.NewToken(kind)
		params := token.Params{
			Config: cfg,
			Name:   test.Name,
			SSH:    ssh.NewToken(test.FS, cfg.SSH),
		}
		tkn := token.NewToken(params)

		Convey("Given I generate a token", t, func() {
			_, gen, err := tkn.Generate(t.Context())
			So(err, ShouldBeNil)

			Convey("When I try to verify", func() {
				_, err := tkn.Verify(t.Context(), gen)

				Convey("Then I should have no error", func() {
					So(err, ShouldBeNil)
				})
			})
		})
	}
}

func TestWithNoConfig(t *testing.T) {
	Convey("When I try to create a token with no config", t, func() {
		params := token.Params{Name: test.Name}
		tkn := token.NewToken(params)

		Convey("Then I should have bo token", func() {
			So(tkn, ShouldBeNil)
		})
	})
}
