package os_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v3/internal/test"
	"github.com/alexfalkowski/go-service/v3/io"
	"github.com/alexfalkowski/go-service/v3/os"
	"github.com/stretchr/testify/require"
)

func TestOpenRootAllowsFileOperationsWithinDirectory(t *testing.T) {
	dir := t.TempDir()
	root, err := os.OpenRoot(dir)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, root.Close())
	})

	file, err := root.OpenFile("file.txt", os.OpenWriteOnly|os.OpenCreate, 0o600)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, file.Close())
	})
	_, err = file.WriteString("contents")
	require.NoError(t, err)

	data, err := test.FS.ReadFile(test.FS.Join(dir, "file.txt"))
	require.NoError(t, err)
	require.Equal(t, "contents", string(data))
}

func TestOpenRootReportsMissingDirectory(t *testing.T) {
	dir := test.FS.Join(t.TempDir(), "missing")

	_, err := os.OpenRoot(dir)
	require.ErrorIs(t, err, io.ErrNotExist)
	var pathErr *os.PathError
	require.ErrorAs(t, err, &pathErr)
	require.Equal(t, dir, pathErr.Path)
}

func TestOpenInRootOpensFileForReading(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, test.FS.WriteFile(test.FS.Join(dir, "file.txt"), []byte("contents"), 0o600))

	file, err := os.OpenInRoot(dir, "file.txt")
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, file.Close())
	})

	data, _, err := io.ReadAll(file)
	require.NoError(t, err)
	require.Equal(t, "contents", string(data))
}

func TestOpenInRootRejectsFileOutsideDirectory(t *testing.T) {
	parent := t.TempDir()
	dir := test.FS.Join(parent, "root")
	require.NoError(t, test.FS.Mkdir(dir, 0o700))
	require.NoError(t, test.FS.WriteFile(test.FS.Join(parent, "file.txt"), []byte("outside"), 0o600))

	name := test.FS.Join("..", "file.txt")
	_, err := os.OpenInRoot(dir, name)
	var pathErr *os.PathError
	require.ErrorAs(t, err, &pathErr)
	require.Equal(t, name, pathErr.Path)
}

//nolint:usetesting
func TestEnvironmentExposesSetAndUnsetVariables(t *testing.T) {
	key := "__ENV_KEY"

	require.NoError(t, os.Setenv(key, "test"))
	require.Equal(t, "test", os.Getenv(key))
	value, ok := os.LookupEnv(key)
	require.True(t, ok)
	require.Equal(t, "test", value)
	require.Contains(t, os.Environ(), key+"=test")
	require.NoError(t, os.Unsetenv(key))
	value, ok = os.LookupEnv(key)
	require.False(t, ok)
	require.Empty(t, value)
}

func TestSanitizeArgsRemovesGoTestArguments(t *testing.T) {
	args := []string{"service", "-test.v", "server", "-config", "config.yml", "-test.run=TestName", "-test-mode"}
	sanitized := os.SanitizeArgs(args)

	require.Equal(t, []string{"service", "server", "-config", "config.yml", "-test-mode"}, sanitized)
	require.Equal(t, []string{"service", "-test.v", "server", "-config", "config.yml", "-test.run=TestName", "-test-mode"}, args)
}
