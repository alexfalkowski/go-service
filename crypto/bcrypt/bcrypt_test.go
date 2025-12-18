package bcrypt_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/bcrypt"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestSigner(t *testing.T) {
	signer := bcrypt.NewSigner()

	s, err := signer.Sign(strings.Bytes("test"))
	require.NoError(t, err)
	require.NotEmpty(t, s)
	require.NoError(t, signer.Verify(s, strings.Bytes("test")))

	s, err = signer.Sign(strings.Bytes("steve"))
	require.NoError(t, err)
	require.Error(t, signer.Verify(s, strings.Bytes("bob")))

	require.Error(t, signer.Verify(strings.Bytes("steve"), strings.Bytes("bob")))
}
