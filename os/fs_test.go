package os_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/stretchr/testify/require"
)

func TestReadFile(t *testing.T) {
	for _, path := range []string{"none"} {
		t.Run(path, func(t *testing.T) {
			_, err := test.FS.ReadFile(path)
			require.Error(t, err)

			require.True(t, test.FS.IsNotExist(err))
			require.False(t, test.FS.PathExists(path))
		})
	}
}

func TestPathExtension(t *testing.T) {
	for _, f := range []string{"file.yaml", "file.test.yaml", "test/.config/existing.client.yaml"} {
		t.Run(f, func(t *testing.T) {
			require.Equal(t, "yaml", test.FS.PathExtension(f))
		})
	}

	require.Empty(t, test.FS.PathExtension("file"))
}

func TestReadSource(t *testing.T) {
	t.Setenv("DUMMY", "yes")
	t.Setenv("EMPTY", "")

	values := []*test.KeyValue[string, string]{
		{Key: "env:DUMMY", Value: "yes"},
		{Key: "env:EMPTY", Value: ""},
		{Key: test.FilePath("configs/invalid.yml"), Value: "not:\n  our:\n    config: test"},
		{Key: "none", Value: "none"},
	}

	for _, value := range values {
		t.Run(value.Key, func(t *testing.T) {
			data, err := test.FS.ReadSource(value.Key)
			require.NoError(t, err)
			require.Equal(t, value.Value, bytes.String(data))
		})
	}
}

func TestReadSourceMissingEnv(t *testing.T) {
	const key = "MISSING_SOURCE"

	require.NoError(t, os.Unsetenv(key))

	_, err := test.FS.ReadSource("env:" + key)
	require.ErrorIs(t, err, os.ErrEnvSourceMissing)
	require.ErrorContains(t, err, "env:"+key)
}

func TestReadSourceMissingEnvName(t *testing.T) {
	_, err := test.FS.ReadSource("env:")
	require.ErrorIs(t, err, os.ErrEnvSourceMissing)
}

func TestPathExistsUsesCleanPath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	path := test.FS.Join(home, "path-exists.txt")
	require.NoError(t, test.FS.WriteFile(path, []byte("ok"), 0o600))

	require.True(t, test.FS.PathExists("~/path-exists.txt"))
}
