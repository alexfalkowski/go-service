package os_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestReadFile(t *testing.T) {
	for _, path := range []string{"none"} {
		_, err := test.FS.ReadFile(path)
		require.Error(t, err)

		require.True(t, test.FS.IsNotExist(err))
		require.False(t, test.FS.PathExists(path))
	}
}

func TestPathExtension(t *testing.T) {
	for _, f := range []string{"file.yaml", "file.test.yaml", "test/.config/existing.client.yaml"} {
		require.Equal(t, "yaml", test.FS.PathExtension(f))
	}

	require.Empty(t, test.FS.PathExtension("file"))
}

func TestReadSource(t *testing.T) {
	t.Setenv("DUMMY", "yes")

	values := []*test.KeyValue[string, string]{
		{Key: "env:DUMMY", Value: "yes"},
		{Key: test.FilePath("configs/invalid.yml"), Value: "not:\n  our:\n    config: test"},
		{Key: "none", Value: "none"},
	}

	for _, value := range values {
		data, err := test.FS.ReadSource(value.Key)
		require.NoError(t, err)
		require.Equal(t, value.Value, bytes.String(data))
	}
}
