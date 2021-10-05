package auth0_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"os"
	"testing"

	"github.com/alexfalkowski/go-service/pkg/cache/ristretto"
	"github.com/alexfalkowski/go-service/pkg/logger/zap"
	"github.com/alexfalkowski/go-service/pkg/security/auth0"
	"github.com/alexfalkowski/go-service/pkg/transport/http"
	"github.com/form3tech-oss/jwt-go"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

const (
	algorithm = "ES256"
)

func TestCorruptToken(t *testing.T) {
	Convey("Given I have a corrupt token", t, func() {
		cfg := &ristretto.Config{
			Name:        "test",
			NumCounters: 1e7,
			MaxCost:     1 << 30,
			BufferItems: 64,
		}

		lc := fxtest.NewLifecycle(t)

		acfg := &auth0.Config{
			URL:           os.Getenv("AUTH0_URL"),
			ClientID:      os.Getenv("AUTH0_CLIENT_ID"),
			ClientSecret:  os.Getenv("AUTH0_CLIENT_SECRET"),
			Audience:      os.Getenv("AUTH0_AUDIENCE"),
			Issuer:        os.Getenv("AUTH0_ISSUER"),
			Algorithm:     os.Getenv("AUTH0_ALGORITHM"),
			JSONWebKeySet: os.Getenv("AUTH0_JSON_WEB_KEY_SET"),
		}

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cache, err := ristretto.NewCache(lc, cfg)
		So(err, ShouldBeNil)

		client := http.NewClient(logger)
		cert := auth0.NewCertificator(acfg, client, cache)
		ver := auth0.NewVerifier(acfg, cert)

		lc.RequireStart()

		Convey("When I try to verify the token", func() {
			ctx := context.Background()
			_, err = ver.Verify(ctx, []byte("corrupt-token"))

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})

		lc.RequireStop()
	})
}

func TestMissingAudienceToken(t *testing.T) {
	Convey("Given I have a missing audience in token", t, func() {
		cfg := &ristretto.Config{
			Name:        "test",
			NumCounters: 1e7,
			MaxCost:     1 << 30,
			BufferItems: 64,
		}

		lc := fxtest.NewLifecycle(t)

		acfg := &auth0.Config{
			URL:           os.Getenv("AUTH0_URL"),
			ClientID:      os.Getenv("AUTH0_CLIENT_ID"),
			ClientSecret:  os.Getenv("AUTH0_CLIENT_SECRET"),
			Audience:      os.Getenv("AUTH0_AUDIENCE"),
			Issuer:        os.Getenv("AUTH0_ISSUER"),
			Algorithm:     os.Getenv("AUTH0_ALGORITHM"),
			JSONWebKeySet: os.Getenv("AUTH0_JSON_WEB_KEY_SET"),
		}

		acfg.Algorithm = algorithm

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cache, err := ristretto.NewCache(lc, cfg)
		So(err, ShouldBeNil)

		client := http.NewClient(logger)
		cert := auth0.NewCertificator(acfg, client, cache)
		ver := auth0.NewVerifier(acfg, cert)

		lc.RequireStart()

		Convey("When I try to verify the token", func() {
			token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{})

			key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			So(err, ShouldBeNil)

			tkn, err := token.SignedString(key)
			So(err, ShouldBeNil)

			ctx := context.Background()
			_, err = ver.Verify(ctx, []byte(tkn))

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})

		lc.RequireStop()
	})
}

func TestMissingIssuerToken(t *testing.T) {
	Convey("Given I have a missing issuer in token", t, func() {
		cfg := &ristretto.Config{
			Name:        "test",
			NumCounters: 1e7,
			MaxCost:     1 << 30,
			BufferItems: 64,
		}

		lc := fxtest.NewLifecycle(t)

		acfg := &auth0.Config{
			URL:           os.Getenv("AUTH0_URL"),
			ClientID:      os.Getenv("AUTH0_CLIENT_ID"),
			ClientSecret:  os.Getenv("AUTH0_CLIENT_SECRET"),
			Audience:      os.Getenv("AUTH0_AUDIENCE"),
			Issuer:        os.Getenv("AUTH0_ISSUER"),
			Algorithm:     os.Getenv("AUTH0_ALGORITHM"),
			JSONWebKeySet: os.Getenv("AUTH0_JSON_WEB_KEY_SET"),
		}

		acfg.Algorithm = algorithm

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cache, err := ristretto.NewCache(lc, cfg)
		So(err, ShouldBeNil)

		client := http.NewClient(logger)
		cert := auth0.NewCertificator(acfg, client, cache)
		ver := auth0.NewVerifier(acfg, cert)

		lc.RequireStart()

		Convey("When I try to verify the token", func() {
			claims := jwt.MapClaims{
				"aud": acfg.Audience,
			}
			token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

			key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			So(err, ShouldBeNil)

			tkn, err := token.SignedString(key)
			So(err, ShouldBeNil)

			ctx := context.Background()
			_, err = ver.Verify(ctx, []byte(tkn))

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})

		lc.RequireStop()
	})
}

func TestInvalidCertificateToken(t *testing.T) {
	Convey("Given I have an invalid jwks endpoint", t, func() {
		cfg := &ristretto.Config{
			Name:        "test",
			NumCounters: 1e7,
			MaxCost:     1 << 30,
			BufferItems: 64,
		}

		lc := fxtest.NewLifecycle(t)

		acfg := &auth0.Config{
			URL:           os.Getenv("AUTH0_URL"),
			ClientID:      os.Getenv("AUTH0_CLIENT_ID"),
			ClientSecret:  os.Getenv("AUTH0_CLIENT_SECRET"),
			Audience:      os.Getenv("AUTH0_AUDIENCE"),
			Issuer:        os.Getenv("AUTH0_ISSUER"),
			Algorithm:     os.Getenv("AUTH0_ALGORITHM"),
			JSONWebKeySet: os.Getenv("AUTH0_JSON_WEB_KEY_SET"),
		}

		acfg.Algorithm = algorithm
		acfg.JSONWebKeySet = "https://non-existent.com/.well-known/jwks.json"

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cache, err := ristretto.NewCache(lc, cfg)
		So(err, ShouldBeNil)

		client := http.NewClient(logger)
		cert := auth0.NewCertificator(acfg, client, cache)
		ver := auth0.NewVerifier(acfg, cert)

		lc.RequireStart()

		Convey("When I try to verify the token", func() {
			claims := jwt.MapClaims{
				"aud": acfg.Audience,
				"iss": acfg.Issuer,
			}
			token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

			key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			So(err, ShouldBeNil)

			tkn, err := token.SignedString(key)
			So(err, ShouldBeNil)

			ctx := context.Background()
			_, err = ver.Verify(ctx, []byte(tkn))

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})

		lc.RequireStop()
	})
}

func TestMissingKidToken(t *testing.T) {
	Convey("Given I have an invalid jwks endpoint", t, func() {
		cfg := &ristretto.Config{
			Name:        "test",
			NumCounters: 1e7,
			MaxCost:     1 << 30,
			BufferItems: 64,
		}

		lc := fxtest.NewLifecycle(t)

		acfg := &auth0.Config{
			URL:           os.Getenv("AUTH0_URL"),
			ClientID:      os.Getenv("AUTH0_CLIENT_ID"),
			ClientSecret:  os.Getenv("AUTH0_CLIENT_SECRET"),
			Audience:      os.Getenv("AUTH0_AUDIENCE"),
			Issuer:        os.Getenv("AUTH0_ISSUER"),
			Algorithm:     os.Getenv("AUTH0_ALGORITHM"),
			JSONWebKeySet: os.Getenv("AUTH0_JSON_WEB_KEY_SET"),
		}

		acfg.Algorithm = algorithm

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cache, err := ristretto.NewCache(lc, cfg)
		So(err, ShouldBeNil)

		client := http.NewClient(logger)
		cert := auth0.NewCertificator(acfg, client, cache)
		ver := auth0.NewVerifier(acfg, cert)

		lc.RequireStart()

		Convey("When I try to verify the token", func() {
			claims := jwt.MapClaims{
				"aud": acfg.Audience,
				"iss": acfg.Issuer,
				"kid": "none",
			}
			token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

			key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			So(err, ShouldBeNil)

			tkn, err := token.SignedString(key)
			So(err, ShouldBeNil)

			ctx := context.Background()
			_, err = ver.Verify(ctx, []byte(tkn))

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})

		lc.RequireStop()
	})
}
