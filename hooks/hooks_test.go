package hooks_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/hooks"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestHooks(t *testing.T) {
	gen := hooks.NewGenerator(rand.NewGenerator(rand.NewReader()))

	c, err := gen.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, c)

	_, err = hooks.NewHook(test.FS, &hooks.Config{Secret: test.FilePath("secrets/none")})
	require.Error(t, err)

	_, err = hooks.NewHook(test.FS, &hooks.Config{Secret: test.FilePath("secrets/redis")})
	require.Error(t, err)

	h, err := hooks.NewHook(nil, nil)
	require.NoError(t, err)
	require.Nil(t, h)
}
