package hooks_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/hooks"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestGenerator(t *testing.T) {
	gen := hooks.NewGenerator(rand.NewGenerator(rand.NewReader()))

	secret, err := gen.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, secret)
}

func TestNewHookReturnsSourceError(t *testing.T) {
	_, err := hooks.NewHook(test.FS, &hooks.Config{Secret: test.FilePath("secrets/none")})
	require.Error(t, err)
}

func TestNewHookReturnsInvalidSecretError(t *testing.T) {
	_, err := hooks.NewHook(test.FS, &hooks.Config{Secret: test.FilePath("secrets/redis")})
	require.Error(t, err)
}

func TestNewHookRejectsEmptySecret(t *testing.T) {
	_, err := hooks.NewHook(test.FS, &hooks.Config{Secret: ""})
	require.ErrorIs(t, err, hooks.ErrEmptySecret)
}

func TestNewHookRejectsEmptyEnvSecret(t *testing.T) {
	t.Setenv("EMPTY_WEBHOOK_VALUE", "")

	source := "env:" + "EMPTY_WEBHOOK_VALUE"
	_, err := hooks.NewHook(test.FS, &hooks.Config{Secret: source})
	require.ErrorIs(t, err, hooks.ErrEmptySecret)
}

func TestNewHookReturnsNilWhenDisabled(t *testing.T) {
	h, err := hooks.NewHook(nil, nil)
	require.NoError(t, err)
	require.Nil(t, h)
}
