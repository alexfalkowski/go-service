package header_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/telemetry/header"
	"github.com/avfs/avfs"
	"github.com/avfs/avfs/vfs/failfs"
	"github.com/avfs/avfs/vfs/memfs"
	"github.com/stretchr/testify/require"
)

func TestSecrets(t *testing.T) {
	headers := header.Map{
		"file":    test.FilePath("secrets/hooks"),
		"literal": "literal-token",
	}

	require.NoError(t, headers.Secrets(test.FS))
	require.Equal(t, header.Map{
		"file":    "QW4xbEphNWtxa09TMWdUY1MydmJybHlKR04zTG5aSEU=",
		"literal": "literal-token",
	}, headers)
	require.Error(t, header.Map{"test": test.FilePath("none")}.Secrets(test.ErrFS))
}

func TestMustSecrets(t *testing.T) {
	headers := header.Map{"test": test.FilePath("secrets/hooks")}

	require.NotPanics(t, func() {
		headers.MustSecrets(test.FS)
	})
	require.Equal(t, header.Map{"test": "QW4xbEphNWtxa09TMWdUY1MydmJybHlKR04zTG5aSEU="}, headers)

	require.Panics(t, func() {
		header.Map{"test": test.FilePath("none")}.MustSecrets(test.ErrFS)
	})
}

func TestSecretsDoNotPartiallyMutateOnError(t *testing.T) {
	base := memfs.New()
	require.NoError(t, base.WriteFile("/first", []byte("one"), 0o644))
	require.NoError(t, base.WriteFile("/second", []byte("two"), 0o644))

	fail := failfs.New(base)
	reads := 0
	require.NoError(t, fail.SetFailFunc(func(_ avfs.VFSBase, fn avfs.FnVFS, _ *failfs.FailParam) error {
		if fn == avfs.FnOpenFile {
			reads++
			if reads > 1 {
				return test.ErrFailed
			}
		}

		return nil
	}))

	fs := &os.FS{VFS: fail}
	headers := header.Map{
		"first":  "file:/first",
		"second": "file:/second",
	}
	original := header.Map{
		"first":  headers["first"],
		"second": headers["second"],
	}

	err := headers.Secrets(fs)
	require.ErrorIs(t, err, test.ErrFailed)
	require.Equal(t, original, headers)
}
