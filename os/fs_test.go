package os_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/stretchr/testify/require"
)

func TestReadFileClassifiesMissingFile(t *testing.T) {
	path := "none"

	_, err := test.FS.ReadFile(path)
	require.Error(t, err)

	require.True(t, test.FS.IsNotExist(err))
	require.False(t, test.FS.PathExists(path))
}

func TestPathExtensionExtractsFinalExtension(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{name: "extracts single extension", path: "file.yaml"},
		{name: "extracts extension after multiple periods", path: "file.test.yaml"},
		{name: "extracts extension from nested path", path: "test/.config/existing.client.yaml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, "yaml", test.FS.PathExtension(tt.path))
		})
	}

	require.Empty(t, test.FS.PathExtension("file"))
}

func TestExpandPathExpandsHomeDirectory(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	require.Empty(t, test.FS.ExpandPath(""))
	require.Equal(t, "path/file.txt", test.FS.ExpandPath("path/file.txt"))
	require.Equal(t, test.FS.Join(home, "path/file.txt"), test.FS.ExpandPath("~/path/file.txt"))
}

func TestCleanPathExpandsBeforeCleaning(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	require.Equal(t, test.FS.Join(home, "..", "path.txt"), test.FS.CleanPath("~/../path.txt"))
}

func TestReadSourceResolvesEnvironmentFileAndLiteralSources(t *testing.T) {
	t.Setenv("DUMMY", "yes")
	t.Setenv("EMPTY", "")

	values := []*test.KeyValue[string, string]{
		{Key: "env:DUMMY", Value: "yes"},
		{Key: "env:EMPTY", Value: ""},
		{Key: test.FilePath("configs/invalid.yaml"), Value: "not:\n  our:\n    config: test"},
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
