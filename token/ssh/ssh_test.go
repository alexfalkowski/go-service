package ssh_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/encoding/base64"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/token/errors"
	"github.com/alexfalkowski/go-service/v2/token/ssh"
	"github.com/stretchr/testify/require"
)

func TestValid(t *testing.T) {
	token := ssh.NewToken(test.NewToken("ssh").SSH, test.FS)

	tkn, err := token.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	sub, err := token.Verify(tkn)
	require.NoError(t, err)
	require.Equal(t, test.UserID.String(), sub)

	ssh := ssh.NewToken(nil, nil)
	require.Nil(t, ssh)
}

func TestValidNameWithDash(t *testing.T) {
	cfg := test.NewToken("ssh").SSH
	cfg.Key.Name = "test-user"
	cfg.Keys[0].Name = "test-user"

	token := ssh.NewToken(cfg, test.FS)

	tkn, err := token.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	sub, err := token.Verify(tkn)
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
	_, err := token.Generate()
	require.Error(t, err)

	for _, tkn := range []string{strings.Empty, "none-", "test-", "test-bob"} {
		t.Run(tkn, func(t *testing.T) {
			token := ssh.NewToken(test.NewToken("ssh").SSH, test.FS)
			sub, err := token.Verify(tkn)
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
	sub, err := token.Verify("test-bob")
	require.Error(t, err)
	require.Empty(t, sub)

	valid := ssh.NewToken(test.NewToken("ssh").SSH, test.FS)
	tkn, err := valid.Generate()
	require.NoError(t, err)

	index := strings.LastIndex(tkn, "-")
	require.NotEqual(t, -1, index)

	sub, err = valid.Verify(tkn[:index+1] + base64.Encode([]byte("bad")))
	require.Error(t, err)
	require.Empty(t, sub)

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

		tkn, err := token.Generate()
		require.Empty(t, tkn)
		require.ErrorIs(t, err, errors.ErrInvalidConfig)
	})

	t.Run("verify with matching key missing config", func(t *testing.T) {
		token := ssh.NewToken(&ssh.Config{
			Keys: ssh.Keys{
				&ssh.Key{Name: "test"},
			},
		}, test.FS)

		sub, err := token.Verify("test-bob")
		require.Empty(t, sub)
		require.ErrorIs(t, err, errors.ErrInvalidConfig)
	})
}
