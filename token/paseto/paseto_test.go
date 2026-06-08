package paseto_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/token/errors"
	"github.com/alexfalkowski/go-service/v2/token/keys"
	"github.com/alexfalkowski/go-service/v2/token/paseto"
	"github.com/stretchr/testify/require"
)

func TestConfigRejectsInvalidValues(t *testing.T) {
	valid := test.NewToken("paseto").Paseto
	tests := []struct {
		config *paseto.Config
		name   string
	}{
		{
			name: "missing issuer",
			config: &paseto.Config{
				Key:        valid.Key,
				Keys:       valid.Keys,
				Expiration: time.Hour,
			},
		},
		{
			name: "missing key",
			config: &paseto.Config{
				Issuer:     "iss",
				Keys:       valid.Keys,
				Expiration: time.Hour,
			},
		},
		{
			name: "missing keys",
			config: &paseto.Config{
				Issuer:     "iss",
				Key:        valid.Key,
				Expiration: time.Hour,
			},
		},
		{
			name: "negative expiration",
			config: &paseto.Config{
				Issuer:     "iss",
				Key:        valid.Key,
				Keys:       valid.Keys,
				Expiration: -time.Second,
			},
		},
		{
			name: "zero expiration",
			config: &paseto.Config{
				Issuer: "iss",
				Key:    valid.Key,
				Keys:   valid.Keys,
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
	cfg := test.NewToken("paseto")
	token := paseto.NewToken(cfg.Paseto, test.FS, uuid.NewGenerator())

	tkn, err := token.Generate("hello", test.UserID.String())
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	sub, err := token.Verify(tkn, "hello")
	require.NoError(t, err)
	require.Equal(t, test.UserID.String(), sub)
}

func TestInvalidEmptySubject(t *testing.T) {
	cfg := test.NewToken("paseto")
	token := paseto.NewToken(cfg.Paseto, test.FS, uuid.NewGenerator())

	tkn, err := token.Generate("hello", strings.Empty)
	require.NoError(t, err)

	sub, err := token.Verify(tkn, "hello")

	require.Empty(t, sub)
	require.ErrorIs(t, err, errors.ErrInvalidSubject)
}

func TestInvalidKeyID(t *testing.T) {
	cfg := test.NewToken("paseto")
	gen := uuid.NewGenerator()
	token := paseto.NewToken(cfg.Paseto, test.FS, gen)

	wrongCfg := cloneConfig(cfg.Paseto)
	wrongCfg.Key = "test"
	wrongCfg.Keys = keys.Map{
		"test": cfg.Paseto.Keys.Get(cfg.Paseto.Key),
	}
	wrong := paseto.NewToken(wrongCfg, test.FS, gen)

	tkn, err := wrong.Generate("hello", test.UserID.String())
	require.NoError(t, err)

	sub, err := token.Verify(tkn, "hello")
	require.Empty(t, sub)
	require.ErrorIs(t, err, errors.ErrInvalidKeyID)
}

func TestInvalidLifetimeExceedsConfig(t *testing.T) {
	cfg := test.NewToken("paseto")
	gen := uuid.NewGenerator()
	generatorCfg := cloneConfig(cfg.Paseto)
	generatorCfg.Expiration = time.Hour
	verifierCfg := cloneConfig(cfg.Paseto)
	verifierCfg.Expiration = time.Minute
	generator := paseto.NewToken(generatorCfg, test.FS, gen)
	verifierToken := paseto.NewToken(verifierCfg, test.FS, gen)

	tkn, err := generator.Generate("hello", test.UserID.String())
	require.NoError(t, err)

	sub, err := verifierToken.Verify(tkn, "hello")
	require.Empty(t, sub)
	require.ErrorIs(t, err, errors.ErrInvalidTime)
}

func TestInvalidIssuer(t *testing.T) {
	cfg := test.NewToken("paseto")
	gen := uuid.NewGenerator()
	generatorCfg := cloneConfig(cfg.Paseto)
	generatorCfg.Issuer = "other"
	generator := paseto.NewToken(generatorCfg, test.FS, gen)
	verifierToken := paseto.NewToken(cfg.Paseto, test.FS, gen)

	tkn, err := generator.Generate("hello", test.UserID.String())
	require.NoError(t, err)

	sub, err := verifierToken.Verify(tkn, "hello")
	require.Empty(t, sub)
	require.Error(t, err)
}

func TestInvalidVerifyExpirationConfig(t *testing.T) {
	cfg := test.NewToken("paseto")
	gen := uuid.NewGenerator()
	generator := paseto.NewToken(cfg.Paseto, test.FS, gen)
	verifierCfg := cloneConfig(cfg.Paseto)
	verifierCfg.Expiration = 0
	verifierToken := paseto.NewToken(verifierCfg, test.FS, gen)

	tkn, err := generator.Generate("hello", test.UserID.String())
	require.NoError(t, err)

	sub, err := verifierToken.Verify(tkn, "hello")
	require.Empty(t, sub)
	require.ErrorIs(t, err, errors.ErrInvalidConfig)
}

func TestInvalidGenerateExpirationConfig(t *testing.T) {
	cfg := test.NewToken("paseto")
	generateCfg := cloneConfig(cfg.Paseto)
	generateCfg.Expiration = 0
	token := paseto.NewToken(generateCfg, test.FS, uuid.NewGenerator())

	tkn, err := token.Generate("hello", test.UserID.String())
	require.Empty(t, tkn)
	require.ErrorIs(t, err, errors.ErrInvalidConfig)
}

func TestInvalidGenerateIssuerConfig(t *testing.T) {
	cfg := test.NewToken("paseto")
	generateCfg := cloneConfig(cfg.Paseto)
	generateCfg.Issuer = strings.Empty
	token := paseto.NewToken(generateCfg, test.FS, uuid.NewGenerator())

	tkn, err := token.Generate("hello", test.UserID.String())
	require.Empty(t, tkn)
	require.ErrorIs(t, err, errors.ErrInvalidConfig)
}

func TestInvalidVerifyIssuerConfig(t *testing.T) {
	cfg := test.NewToken("paseto")
	gen := uuid.NewGenerator()
	generator := paseto.NewToken(cfg.Paseto, test.FS, gen)
	verifierCfg := cloneConfig(cfg.Paseto)
	verifierCfg.Issuer = strings.Empty
	verifierToken := paseto.NewToken(verifierCfg, test.FS, gen)

	tkn, err := generator.Generate("hello", test.UserID.String())
	require.NoError(t, err)

	sub, err := verifierToken.Verify(tkn, "hello")
	require.Empty(t, sub)
	require.ErrorIs(t, err, errors.ErrInvalidConfig)
}

func TestInvalidAudience(t *testing.T) {
	gen := uuid.NewGenerator()
	cfg := test.NewToken("paseto")
	token := paseto.NewToken(cfg.Paseto, test.FS, gen)

	tkn, err := token.Generate("hello", test.UserID.String())
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	_, err = token.Verify(tkn, "test")
	require.Error(t, err)
}

func TestValidMatchingIssuerAndAudience(t *testing.T) {
	cfg := test.NewToken("paseto")
	cfg.Paseto.Issuer = "test"
	token := paseto.NewToken(cfg.Paseto, test.FS, uuid.NewGenerator())

	tkn, err := token.Generate("hello", test.UserID.String())
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	_, err = token.Verify(tkn, "hello")
	require.NoError(t, err)
}

func TestInvalidMalformedToken(t *testing.T) {
	cfg := test.NewToken("paseto")
	token := paseto.NewToken(cfg.Paseto, test.FS, uuid.NewGenerator())

	_, err := token.Verify("invalid", "aud")
	require.Error(t, err)
}

func TestInvalidGenerateMalformedPrivateKey(t *testing.T) {
	cfg := test.NewToken("paseto")
	badPrivate := cloneConfig(cfg.Paseto)
	badPrivate.Keys = keys.Map{
		badPrivate.Key: &keys.Config{Config: &ed25519.Config{
			Public:  test.NewEd25519().Public,
			Private: "short",
		}},
	}
	token := paseto.NewToken(badPrivate, test.FS, uuid.NewGenerator())

	_, err := token.Generate("hello", test.UserID.String())
	require.Error(t, err)
}

func TestInvalidVerifyMalformedPublicKey(t *testing.T) {
	cfg := test.NewToken("paseto")
	token := paseto.NewToken(cfg.Paseto, test.FS, uuid.NewGenerator())

	_, err := token.Verify(strings.Empty, "aud")
	require.Error(t, err)
}

func TestInvalidNilConfig(t *testing.T) {
	token := paseto.NewToken(nil, test.FS, uuid.NewGenerator())
	require.Nil(t, token)
}

func TestInvalidConfigDoesNotPanic(t *testing.T) {
	cfg := test.NewToken("paseto")
	gen := uuid.NewGenerator()

	t.Run("generate without private key", func(t *testing.T) {
		noPrivate := cloneConfig(cfg.Paseto)
		noPrivate.Keys = keys.Map{
			noPrivate.Key: &keys.Config{Config: &ed25519.Config{Public: test.NewEd25519().Public}},
		}
		token := paseto.NewToken(noPrivate, test.FS, gen)

		tkn, err := token.Generate("hello", test.UserID.String())
		require.Empty(t, tkn)
		require.Error(t, err)
	})

	t.Run("generate without generator", func(t *testing.T) {
		token := paseto.NewToken(cfg.Paseto, test.FS, nil)

		tkn, err := token.Generate("hello", test.UserID.String())
		require.Empty(t, tkn)
		require.ErrorIs(t, err, errors.ErrInvalidConfig)
	})

	t.Run("verify without public key", func(t *testing.T) {
		valid := paseto.NewToken(cfg.Paseto, test.FS, gen)
		tkn, err := valid.Generate("hello", test.UserID.String())
		require.NoError(t, err)

		noPublic := cloneConfig(cfg.Paseto)
		noPublic.Keys = keys.Map{
			noPublic.Key: &keys.Config{Config: &ed25519.Config{Private: test.NewEd25519().Private}},
		}
		token := paseto.NewToken(noPublic, test.FS, gen)

		sub, err := token.Verify(tkn, "hello")
		require.Empty(t, sub)
		require.Error(t, err)
	})
}

func TestInvalidKeyMaterialDoesNotPanic(t *testing.T) {
	cfg := test.NewToken("paseto")
	gen := uuid.NewGenerator()

	t.Run("generate with malformed private key", func(t *testing.T) {
		badPrivate := cloneConfig(cfg.Paseto)
		badPrivate.Keys = keys.Map{
			badPrivate.Key: &keys.Config{Config: &ed25519.Config{
				Public:  test.NewEd25519().Public,
				Private: "short",
			}},
		}
		token := paseto.NewToken(badPrivate, test.FS, gen)

		tkn, err := token.Generate("hello", test.UserID.String())
		require.Empty(t, tkn)
		require.Error(t, err)
	})

	t.Run("verify with malformed public key", func(t *testing.T) {
		valid := paseto.NewToken(cfg.Paseto, test.FS, gen)
		tkn, err := valid.Generate("hello", test.UserID.String())
		require.NoError(t, err)

		badPublic := cloneConfig(cfg.Paseto)
		badPublic.Keys = keys.Map{
			badPublic.Key: &keys.Config{Config: &ed25519.Config{
				Public:  "short",
				Private: test.NewEd25519().Private,
			}},
		}
		token := paseto.NewToken(badPublic, test.FS, gen)

		sub, err := token.Verify(tkn, "hello")
		require.Empty(t, sub)
		require.Error(t, err)
	})
}

func cloneConfig(cfg *paseto.Config) *paseto.Config {
	return &paseto.Config{
		Issuer:     cfg.Issuer,
		Key:        cfg.Key,
		Keys:       cfg.Keys,
		Expiration: cfg.Expiration,
	}
}
