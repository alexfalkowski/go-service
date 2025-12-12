package ssh_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
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
		token := ssh.NewToken(test.NewToken("ssh").SSH, test.FS)
		_, err := token.Verify(tkn)
		require.Error(t, err)
	}

	token = ssh.NewToken(&ssh.Config{
		Keys: ssh.Keys{
			&ssh.Key{
				Name:   "test",
				Config: test.NewSSH("secrets/none", "secrets/ssh_private"),
			},
		},
	}, test.FS)
	_, err = token.Verify("test-bob")
	require.Error(t, err)

	token = ssh.NewToken(nil, test.FS)
	require.Nil(t, token)
}
