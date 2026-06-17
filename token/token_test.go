package token_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/token/errors"
	"github.com/alexfalkowski/go-service/v2/token/jwt"
	"github.com/alexfalkowski/go-service/v2/token/paseto"
	"github.com/alexfalkowski/go-service/v2/token/ssh"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	for _, kind := range []string{"jwt", "paseto"} {
		t.Run(kind, func(t *testing.T) {
			cfg := test.NewToken(kind)
			gen := uuid.NewGenerator()
			tkn := token.NewToken(test.Name, cfg, test.FS, gen)

			_, err := tkn.Generate("hello", test.UserID.String())
			require.NoError(t, err)
		})
	}
}

func TestVerify(t *testing.T) {
	for _, kind := range []string{"jwt", "paseto"} {
		t.Run(kind, func(t *testing.T) {
			cfg := test.NewToken(kind)
			gen := uuid.NewGenerator()
			tkn := token.NewToken(test.Name, cfg, test.FS, gen)

			token, err := tkn.Generate("hello", test.UserID.String())
			require.NoError(t, err)

			sub, err := tkn.Verify(token, "hello")
			require.NoError(t, err)
			require.Equal(t, test.UserID.String(), sub)

			sub, err = tkn.Verify(token, "other")
			require.Equal(t, strings.Empty, sub)
			require.Error(t, err)
			require.NotErrorIs(t, err, errors.ErrInvalidMatch)
		})
	}

	t.Run("ssh", func(t *testing.T) {
		cfg := test.NewToken("ssh")
		tkn := token.NewToken(test.Name, cfg, test.FS, nil)

		gen, err := tkn.Generate("hello", strings.Empty)
		require.NoError(t, err)

		sshToken := ssh.NewToken(cfg.SSH, test.FS)
		sub, err := sshToken.Verify(bytes.String(gen), "hello")
		require.NoError(t, err)
		require.Equal(t, test.UserID.String(), sub)

		sub, err = tkn.Verify(gen, "hello")
		require.NoError(t, err)
		require.Equal(t, test.UserID.String(), sub)

		_, err = tkn.Verify(gen, "other")
		require.ErrorIs(t, err, errors.ErrInvalidAudience)
	})
}

func TestVerifyDispatchesJWTToConcreteImplementation(t *testing.T) {
	cfg := test.NewToken("jwt")
	gen := uuid.NewGenerator()
	tkn := token.NewToken(test.Name, cfg, test.FS, gen)

	raw, err := tkn.Generate("hello", test.UserID.String())
	require.NoError(t, err)

	jwtToken := jwt.NewToken(cfg.JWT, test.FS, gen)
	sub, err := jwtToken.Verify(bytes.String(raw), "hello")
	require.NoError(t, err)
	require.Equal(t, test.UserID.String(), sub)
}

func TestVerifyDispatchesPasetoToConcreteImplementation(t *testing.T) {
	cfg := test.NewToken("paseto")
	gen := uuid.NewGenerator()
	tkn := token.NewToken(test.Name, cfg, test.FS, gen)

	raw, err := tkn.Generate("hello", test.UserID.String())
	require.NoError(t, err)

	pasetoToken := paseto.NewToken(cfg.Paseto, test.FS, gen)
	sub, err := pasetoToken.Verify(bytes.String(raw), "hello")
	require.NoError(t, err)
	require.Equal(t, test.UserID.String(), sub)
}

func TestVerifyDispatchesSSHToConcreteImplementation(t *testing.T) {
	cfg := test.NewToken("ssh")
	tkn := token.NewToken(test.Name, cfg, test.FS, nil)

	raw, err := tkn.Generate("hello", strings.Empty)
	require.NoError(t, err)

	sshToken := ssh.NewToken(cfg.SSH, test.FS)
	sub, err := sshToken.Verify(bytes.String(raw), "hello")
	require.NoError(t, err)
	require.Equal(t, test.UserID.String(), sub)
}

func TestVerifyRejectsPasetoEmptySubject(t *testing.T) {
	cfg := test.NewToken("paseto")
	gen := uuid.NewGenerator()
	tkn := token.NewToken(test.Name, cfg, test.FS, gen)

	raw, err := tkn.Generate("hello", strings.Empty)
	require.NoError(t, err)

	sub, err := tkn.Verify(raw, "hello")
	require.Equal(t, strings.Empty, sub)
	require.ErrorIs(t, err, errors.ErrInvalidSubject)
}

func TestVerifyRejectsJWTEmptySubject(t *testing.T) {
	cfg := test.NewToken("jwt")
	gen := uuid.NewGenerator()
	tkn := token.NewToken(test.Name, cfg, test.FS, gen)

	raw, err := tkn.Generate("hello", strings.Empty)
	require.NoError(t, err)

	sub, err := tkn.Verify(raw, "hello")
	require.Equal(t, strings.Empty, sub)
	require.ErrorIs(t, err, errors.ErrInvalidSubject)
}

func TestSSHSubjectMatchesActiveKey(t *testing.T) {
	cfg := test.NewToken("ssh")
	tkn := token.NewToken(test.Name, cfg, test.FS, nil)

	raw, err := tkn.Generate("hello", "ignored-subject")
	require.NoError(t, err)

	sub, err := tkn.Verify(raw, "hello")
	require.NoError(t, err)
	require.Equal(t, test.UserID.String(), sub)
}

func TestConfigRejectsInvalidValues(t *testing.T) {
	valid := test.NewToken("jwt")
	tests := []struct {
		config *token.Config
		name   string
	}{
		{
			name:   "missing kind",
			config: &token.Config{JWT: valid.JWT},
		},
		{
			name:   "unknown kind",
			config: &token.Config{Kind: "none"},
		},
		{
			name:   "missing jwt config",
			config: &token.Config{Kind: "jwt"},
		},
		{
			name:   "missing paseto config",
			config: &token.Config{Kind: "paseto"},
		},
		{
			name:   "missing ssh config",
			config: &token.Config{Kind: "ssh"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Error(t, test.Validator.Struct(tt.config))
		})
	}

	for _, kind := range []string{"jwt", "paseto", "ssh"} {
		t.Run(kind, func(t *testing.T) {
			require.NoError(t, test.Validator.Struct(test.NewToken(kind)))
		})
	}
}

func TestUnknownKindConfig(t *testing.T) {
	cfg := test.NewToken("none")
	tkn := token.NewToken(test.Name, cfg, test.FS, nil)

	gen, err := tkn.Generate("hello", test.UserID.String())
	require.Nil(t, gen)
	require.ErrorIs(t, err, errors.ErrInvalidConfig)

	sub, err := tkn.Verify([]byte("test"), "hello")
	require.Equal(t, strings.Empty, sub)
	require.ErrorIs(t, err, errors.ErrInvalidConfig)
}

func TestNewTokenWithNilConfig(t *testing.T) {
	tkn := token.NewToken(test.Name, nil, test.FS, nil)
	require.Nil(t, tkn)
}

func TestInvalidKindConfig(t *testing.T) {
	for _, kind := range []string{"jwt", "paseto", "ssh"} {
		t.Run(kind, func(t *testing.T) {
			cfg := &token.Config{Kind: kind}
			tkn := token.NewToken(test.Name, cfg, test.FS, nil)

			gen, err := tkn.Generate("hello", test.UserID.String())
			require.Nil(t, gen)
			require.ErrorIs(t, err, errors.ErrInvalidConfig)

			sub, err := tkn.Verify([]byte("test"), "hello")
			require.Equal(t, strings.Empty, sub)
			require.ErrorIs(t, err, errors.ErrInvalidConfig)
		})
	}
}

func TestInvalidDependenciesDoNotPanic(t *testing.T) {
	for _, kind := range []string{"jwt", "paseto"} {
		t.Run(kind, func(t *testing.T) {
			cfg := test.NewToken(kind)
			tkn := token.NewToken(test.Name, cfg, test.FS, nil)

			gen, err := tkn.Generate("hello", test.UserID.String())
			require.Nil(t, gen)
			require.ErrorIs(t, err, errors.ErrInvalidConfig)
		})
	}
}

func TestInvalidMatchClassification(t *testing.T) {
	gen := uuid.NewGenerator()

	for _, kind := range []string{"jwt", "paseto"} {
		t.Run(kind, func(t *testing.T) {
			cfg := test.NewToken(kind)
			tkn := token.NewToken(test.Name, cfg, test.FS, gen)

			sub, err := tkn.Verify([]byte("test"), "hello")
			require.Equal(t, strings.Empty, sub)
			require.ErrorIs(t, err, errors.ErrInvalidMatch)
		})
	}

	t.Run("ssh", func(t *testing.T) {
		cfg := test.NewToken("ssh")
		tkn := token.NewToken(test.Name, cfg, test.FS, nil)

		sub, err := tkn.Verify([]byte("test"), "hello")
		require.Equal(t, strings.Empty, sub)
		require.ErrorIs(t, err, errors.ErrInvalidMatch)
	})
}
