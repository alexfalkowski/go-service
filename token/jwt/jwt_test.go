package jwt_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/token/errors"
	"github.com/alexfalkowski/go-service/v2/token/jwt"
	"github.com/alexfalkowski/go-service/v2/token/keys"
	"github.com/stretchr/testify/require"
)

func TestConfigRejectsInvalidValues(t *testing.T) {
	valid := test.NewToken("jwt").JWT
	tests := []struct {
		config *jwt.Config
		name   string
	}{
		{
			name: "missing issuer",
			config: &jwt.Config{
				Key:        valid.Key,
				Keys:       valid.Keys,
				Expiration: time.Hour,
			},
		},
		{
			name: "missing key",
			config: &jwt.Config{
				Issuer:     "iss",
				Keys:       valid.Keys,
				Expiration: time.Hour,
			},
		},
		{
			name: "missing keys",
			config: &jwt.Config{
				Issuer:     "iss",
				Key:        valid.Key,
				Expiration: time.Hour,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Error(t, test.Validator.Struct(tt.config))
		})
	}

	require.NoError(t, test.Validator.Struct(valid))
}

func TestValid(t *testing.T) {
	cfg := test.NewToken("jwt")
	token := jwt.NewToken(cfg.JWT, test.FS, uuid.NewGenerator())

	tkn, err := token.Generate("hello", test.UserID.String())
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	sub, err := token.Verify(tkn, "hello")
	require.NoError(t, err)
	require.Equal(t, test.UserID.String(), sub)
}

func TestInvalid(t *testing.T) {
	cfg := test.NewToken("jwt")
	gen := uuid.NewGenerator()
	token := jwt.NewToken(cfg.JWT, test.FS, gen)

	tests := []struct {
		name  string
		value string
	}{
		{
			name:  "malformed token",
			value: "invalid",
		},
		{
			name:  "unexpected signing method",
			value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := token.Verify(tt.value, "hello")
			require.Error(t, err)
		})
	}

	t.Run("invalid audience", func(t *testing.T) {
		tkn, err := token.Generate("hello", test.UserID.String())
		require.NoError(t, err)
		require.NotEmpty(t, tkn)

		sub, err := token.Verify(tkn, "test")
		require.Empty(t, sub)
		require.ErrorIs(t, err, errors.ErrInvalidAudience)
	})

	t.Run("invalid issuer", func(t *testing.T) {
		issuerCfg := cloneConfig(cfg.JWT)
		issuerCfg.Issuer = "test"
		issuerToken := jwt.NewToken(issuerCfg, test.FS, gen)

		tkn, err := issuerToken.Generate("hello", test.UserID.String())
		require.NoError(t, err)
		require.NotEmpty(t, tkn)

		sub, err := token.Verify(tkn, "hello")
		require.Empty(t, sub)
		require.ErrorIs(t, err, errors.ErrInvalidIssuer)
	})

	t.Run("disabled config", func(t *testing.T) {
		token := jwt.NewToken(nil, test.FS, gen)
		require.Nil(t, token)
	})
}

func TestInvalidSignature(t *testing.T) {
	cfg := test.NewToken("jwt")
	token := jwt.NewToken(cfg.JWT, test.FS, uuid.NewGenerator())

	_, private, err := ed25519.NewGenerator(rand.NewGenerator(rand.NewReader())).Generate()
	require.NoError(t, err)

	wrongCfg := cloneConfig(cfg.JWT)
	wrongCfg.Keys = keys.Map{
		wrongCfg.Key: {
			Config: &ed25519.Config{
				Public:  test.NewEd25519().Public,
				Private: private,
			},
		},
	}
	wrong := jwt.NewToken(wrongCfg, test.FS, uuid.NewGenerator())

	tkn, err := wrong.Generate("hello", test.UserID.String())
	require.NoError(t, err)

	sub, err := token.Verify(tkn, "hello")
	require.Empty(t, sub)
	require.ErrorIs(t, err, jwt.ErrTokenSignatureInvalid)
}

func TestInvalidKeyID(t *testing.T) {
	cfg := test.NewToken("jwt")
	gen := uuid.NewGenerator()
	token := jwt.NewToken(cfg.JWT, test.FS, gen)

	t.Run("unexpected key id", func(t *testing.T) {
		wrongCfg := cloneConfig(cfg.JWT)
		wrongCfg.Key = "test"
		wrongCfg.Keys = keys.Map{
			"test": cfg.JWT.Keys.Get(cfg.JWT.Key),
		}
		wrong := jwt.NewToken(wrongCfg, test.FS, gen)

		tkn, err := wrong.Generate("hello", test.UserID.String())
		require.NoError(t, err)
		require.NotEmpty(t, tkn)

		sub, err := token.Verify(tkn, "hello")
		require.Empty(t, sub)
		require.ErrorIs(t, err, errors.ErrInvalidKeyID)
	})

	t.Run("missing key id", func(t *testing.T) {
		tkn := signedJWT(t, cfg.JWT, func(header map[string]any) {
			delete(header, "kid")
		}, nil)

		sub, err := token.Verify(tkn, "hello")
		require.Empty(t, sub)
		require.ErrorIs(t, err, errors.ErrInvalidKeyID)
	})

	t.Run("empty key id", func(t *testing.T) {
		tkn := signedJWT(t, cfg.JWT, func(header map[string]any) {
			header["kid"] = ""
		}, nil)

		sub, err := token.Verify(tkn, "hello")
		require.Empty(t, sub)
		require.ErrorIs(t, err, errors.ErrInvalidKeyID)
	})

	t.Run("non string key id", func(t *testing.T) {
		tkn := signedJWT(t, cfg.JWT, func(header map[string]any) {
			header["kid"] = 123
		}, nil)

		sub, err := token.Verify(tkn, "hello")
		require.Empty(t, sub)
		require.ErrorIs(t, err, errors.ErrInvalidKeyID)
	})
}

func TestInvalidRequiredClaims(t *testing.T) {
	cfg := test.NewToken("jwt")
	token := jwt.NewToken(cfg.JWT, test.FS, uuid.NewGenerator())

	tests := []struct {
		claims func(*jwt.RegisteredClaims)
		name   string
	}{
		{
			name: "missing expiration",
			claims: func(claims *jwt.RegisteredClaims) {
				claims.ExpiresAt = nil
			},
		},
		{
			name: "missing issued at",
			claims: func(claims *jwt.RegisteredClaims) {
				claims.IssuedAt = nil
			},
		},
		{
			name: "missing not before",
			claims: func(claims *jwt.RegisteredClaims) {
				claims.NotBefore = nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tkn := signedJWT(t, cfg.JWT, nil, tt.claims)

			sub, err := token.Verify(tkn, "hello")
			require.Empty(t, sub)
			require.ErrorIs(t, err, errors.ErrInvalidTime)
		})
	}
}

func TestInvalidLifetimeExceedsConfig(t *testing.T) {
	cfg := test.NewToken("jwt")
	gen := uuid.NewGenerator()
	generatorCfg := cloneConfig(cfg.JWT)
	generatorCfg.Expiration = time.Hour
	verifierCfg := cloneConfig(cfg.JWT)
	verifierCfg.Expiration = time.Minute
	generator := jwt.NewToken(generatorCfg, test.FS, gen)
	verifierToken := jwt.NewToken(verifierCfg, test.FS, gen)

	tkn, err := generator.Generate("hello", test.UserID.String())
	require.NoError(t, err)

	sub, err := verifierToken.Verify(tkn, "hello")
	require.Empty(t, sub)
	require.ErrorIs(t, err, errors.ErrInvalidTime)
}

func TestInvalidVerifyExpirationConfig(t *testing.T) {
	cfg := test.NewToken("jwt")
	gen := uuid.NewGenerator()
	generator := jwt.NewToken(cfg.JWT, test.FS, gen)
	verifierCfg := cloneConfig(cfg.JWT)
	verifierCfg.Expiration = 0
	verifierToken := jwt.NewToken(verifierCfg, test.FS, gen)

	tkn, err := generator.Generate("hello", test.UserID.String())
	require.NoError(t, err)

	sub, err := verifierToken.Verify(tkn, "hello")
	require.Empty(t, sub)
	require.ErrorIs(t, err, errors.ErrInvalidConfig)
}

func TestInvalidConfigDoesNotPanic(t *testing.T) {
	cfg := test.NewToken("jwt")
	gen := uuid.NewGenerator()

	t.Run("generate without private key", func(t *testing.T) {
		noPrivate := cloneConfig(cfg.JWT)
		noPrivate.Keys = keys.Map{
			noPrivate.Key: &keys.Config{Config: &ed25519.Config{Public: test.NewEd25519().Public}},
		}
		token := jwt.NewToken(noPrivate, test.FS, gen)

		tkn, err := token.Generate("hello", test.UserID.String())
		require.Empty(t, tkn)
		require.Error(t, err)
	})

	t.Run("generate without generator", func(t *testing.T) {
		token := jwt.NewToken(cfg.JWT, test.FS, nil)

		tkn, err := token.Generate("hello", test.UserID.String())
		require.Empty(t, tkn)
		require.ErrorIs(t, err, errors.ErrInvalidConfig)
	})

	t.Run("verify without public key", func(t *testing.T) {
		valid := jwt.NewToken(cfg.JWT, test.FS, gen)
		tkn, err := valid.Generate("hello", test.UserID.String())
		require.NoError(t, err)

		noPublic := cloneConfig(cfg.JWT)
		noPublic.Keys = keys.Map{
			noPublic.Key: &keys.Config{Config: &ed25519.Config{Private: test.NewEd25519().Private}},
		}
		token := jwt.NewToken(noPublic, test.FS, gen)

		sub, err := token.Verify(tkn, "hello")
		require.Empty(t, sub)
		require.Error(t, err)
	})
}

func TestInvalidGenerateConfigDoesNotPanic(t *testing.T) {
	cfg := test.NewToken("jwt")
	gen := uuid.NewGenerator()

	for _, tt := range []struct {
		config *jwt.Config
		name   string
	}{
		{
			name: "generate without issuer",
			config: &jwt.Config{
				Key:        cfg.JWT.Key,
				Keys:       cfg.JWT.Keys,
				Expiration: cfg.JWT.Expiration,
			},
		},
		{
			name: "generate without key",
			config: &jwt.Config{
				Issuer:     cfg.JWT.Issuer,
				Keys:       cfg.JWT.Keys,
				Expiration: cfg.JWT.Expiration,
			},
		},
		{
			name: "generate without expiration",
			config: &jwt.Config{
				Issuer: cfg.JWT.Issuer,
				Key:    cfg.JWT.Key,
				Keys:   cfg.JWT.Keys,
			},
		},
		{
			name: "generate with negative expiration",
			config: &jwt.Config{
				Issuer:     cfg.JWT.Issuer,
				Key:        cfg.JWT.Key,
				Keys:       cfg.JWT.Keys,
				Expiration: -time.Second,
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			token := jwt.NewToken(tt.config, test.FS, gen)

			tkn, err := token.Generate("hello", test.UserID.String())
			require.Empty(t, tkn)
			require.ErrorIs(t, err, errors.ErrInvalidConfig)
		})
	}
}

func TestInvalidPrivateKeyDoesNotPanic(t *testing.T) {
	cfg := test.NewToken("jwt")
	badPrivate := cloneConfig(cfg.JWT)
	badPrivate.Keys = keys.Map{
		badPrivate.Key: &keys.Config{Config: &ed25519.Config{
			Public:  test.NewEd25519().Public,
			Private: "short",
		}},
	}
	token := jwt.NewToken(badPrivate, test.FS, uuid.NewGenerator())

	tkn, err := token.Generate("hello", test.UserID.String())
	require.Empty(t, tkn)
	require.Error(t, err)
}

func signedJWT(
	t *testing.T,
	cfg *jwt.Config,
	header func(map[string]any),
	claimsFunc func(*jwt.RegisteredClaims),
) string {
	t.Helper()

	signer, err := cfg.Keys.Get(cfg.Key).Signer(test.PEM)
	require.NoError(t, err)

	now := time.Now()
	claims := &jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{Time: now.Add(cfg.Expiration.Duration())},
		ID:        "test-id",
		IssuedAt:  &jwt.NumericDate{Time: now},
		Issuer:    cfg.Issuer,
		NotBefore: &jwt.NumericDate{Time: now},
		Audience:  []string{"hello"},
		Subject:   test.UserID.String(),
	}
	if claimsFunc != nil {
		claimsFunc(claims)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = cfg.Key
	if header != nil {
		header(token.Header)
	}

	tkn, err := token.SignedString(signer.PrivateKey)
	require.NoError(t, err)

	return tkn
}

func cloneConfig(cfg *jwt.Config) *jwt.Config {
	return &jwt.Config{
		Issuer:     cfg.Issuer,
		Key:        cfg.Key,
		Keys:       cfg.Keys,
		Expiration: cfg.Expiration,
	}
}
