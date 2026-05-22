package ssh_test

import (
	"testing"

	crypto "github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/encoding/base64"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/token/errors"
	"github.com/alexfalkowski/go-service/v2/token/ssh"
	"github.com/stretchr/testify/require"
)

func TestConfigRejectsNegativeExpiration(t *testing.T) {
	cfg := &ssh.Config{Expiration: -time.Second}
	require.Error(t, test.Validator.Struct(cfg))
}

func TestValid(t *testing.T) {
	token := ssh.NewToken(test.NewToken("ssh").SSH, test.FS)

	tkn, err := token.Generate(strings.Empty, strings.Empty)
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	sub, err := token.Verify(tkn, strings.Empty)
	require.NoError(t, err)
	require.Equal(t, test.UserID.String(), sub)

	ssh := ssh.NewToken(nil, nil)
	require.Nil(t, ssh)
}

func TestValidForAudience(t *testing.T) {
	token := ssh.NewToken(test.NewToken("ssh").SSH, test.FS)

	tkn, err := token.Generate("/service.Method", strings.Empty)
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	sub, err := token.Verify(tkn, "/service.Method")
	require.NoError(t, err)
	require.Equal(t, test.UserID.String(), sub)
}

func TestInvalidAudience(t *testing.T) {
	token := ssh.NewToken(test.NewToken("ssh").SSH, test.FS)

	tkn, err := token.Generate("/service.Method", strings.Empty)
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	sub, err := token.Verify(tkn, "/service.Other")
	require.Empty(t, sub)
	require.ErrorIs(t, err, errors.ErrInvalidAudience)
}

func TestInvalidExpired(t *testing.T) {
	cfg := test.NewToken("ssh").SSH
	cfg.Expiration = time.Nanosecond
	token := ssh.NewToken(cfg, test.FS)

	tkn, err := token.Generate("/service.Method", strings.Empty)
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	time.Sleep(time.Millisecond)

	sub, err := token.Verify(tkn, "/service.Method")
	require.Empty(t, sub)
	require.ErrorIs(t, err, errors.ErrInvalidTime)
}

func TestValidNameWithDash(t *testing.T) {
	cfg := test.NewToken("ssh").SSH
	cfg.Key.Name = "test-user"
	cfg.Keys[0].Name = "test-user"

	token := ssh.NewToken(cfg, test.FS)

	tkn, err := token.Generate(strings.Empty, strings.Empty)
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	sub, err := token.Verify(tkn, strings.Empty)
	require.NoError(t, err)
	require.Equal(t, "test-user", sub)
}

func TestInvalid(t *testing.T) {
	token := ssh.NewToken(&ssh.Config{
		Key: &ssh.Key{
			Name:   "test",
			Config: test.NewSSH("secrets/ssh_public", "secrets/none"),
		},
	}, test.FS)
	_, err := token.Generate(strings.Empty, strings.Empty)
	require.Error(t, err)

	for _, tkn := range []string{strings.Empty, "none-", "test-", "test-bob"} {
		t.Run(tkn, func(t *testing.T) {
			token := ssh.NewToken(test.NewToken("ssh").SSH, test.FS)
			sub, err := token.Verify(tkn, strings.Empty)
			require.Error(t, err)
			require.Empty(t, sub)
		})
	}

	token = ssh.NewToken(&ssh.Config{
		Keys: ssh.Keys{
			&ssh.Key{
				Name:   "test",
				Config: test.NewSSH("secrets/none", "secrets/ssh_private"),
			},
		},
	}, test.FS)
	sub, err := token.Verify("test-bob", strings.Empty)
	require.Error(t, err)
	require.Empty(t, sub)

	valid := ssh.NewToken(test.NewToken("ssh").SSH, test.FS)
	tkn, err := valid.Generate(strings.Empty, strings.Empty)
	require.NoError(t, err)

	token = ssh.NewToken(&ssh.Config{
		Keys: ssh.Keys{
			&ssh.Key{
				Name:   "other",
				Config: test.NewSSH("secrets/ssh_public", "secrets/ssh_private"),
			},
		},
	}, test.FS)
	sub, err = token.Verify(tkn, strings.Empty)
	require.Empty(t, sub)
	require.ErrorIs(t, err, crypto.ErrInvalidMatch)

	encoded, _, ok := strings.Cut(tkn, ".")
	require.True(t, ok)

	sub, err = valid.Verify(encoded+"."+base64.Encode([]byte("bad")), "other")
	require.Empty(t, sub)
	require.ErrorIs(t, err, crypto.ErrInvalidMatch)

	token = ssh.NewToken(nil, test.FS)
	require.Nil(t, token)
}

func TestInvalidConfigDoesNotPanic(t *testing.T) {
	t.Run("generate with verification only config", func(t *testing.T) {
		token := ssh.NewToken(&ssh.Config{
			Keys: ssh.Keys{
				&ssh.Key{
					Name:   "test",
					Config: test.NewSSH("secrets/ssh_public", "secrets/ssh_private"),
				},
			},
		}, test.FS)

		tkn, err := token.Generate(strings.Empty, strings.Empty)
		require.Empty(t, tkn)
		require.ErrorIs(t, err, errors.ErrInvalidConfig)
	})

	t.Run("generate with empty key name", func(t *testing.T) {
		cfg := test.NewToken("ssh").SSH
		cfg.Key.Name = strings.Empty

		token := ssh.NewToken(cfg, test.FS)

		tkn, err := token.Generate(strings.Empty, strings.Empty)
		require.Empty(t, tkn)
		require.ErrorIs(t, err, errors.ErrInvalidConfig)
	})

	t.Run("generate with empty expiration", func(t *testing.T) {
		cfg := test.NewToken("ssh").SSH
		cfg.Expiration = 0

		token := ssh.NewToken(cfg, test.FS)

		tkn, err := token.Generate(strings.Empty, strings.Empty)
		require.Empty(t, tkn)
		require.ErrorIs(t, err, errors.ErrInvalidConfig)
	})

	t.Run("verify with matching key missing config", func(t *testing.T) {
		valid := ssh.NewToken(test.NewToken("ssh").SSH, test.FS)
		tkn, err := valid.Generate(strings.Empty, strings.Empty)
		require.NoError(t, err)

		token := ssh.NewToken(&ssh.Config{
			Keys: ssh.Keys{
				&ssh.Key{Name: test.UserID.String()},
			},
		}, test.FS)

		sub, err := token.Verify(tkn, strings.Empty)
		require.Empty(t, sub)
		require.ErrorIs(t, err, errors.ErrInvalidConfig)
	})
}

func TestKeysGetIgnoresNilEntries(t *testing.T) {
	keys := ssh.Keys{
		nil,
		&ssh.Key{Name: "other"},
		nil,
		&ssh.Key{Name: "test"},
	}

	key := keys.Get("test")
	require.NotNil(t, key)
	require.Equal(t, "test", key.Name)
	require.Nil(t, ssh.Keys{nil, nil}.Get("missing"))
}
