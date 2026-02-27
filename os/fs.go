package os

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/avfs/avfs"
	"github.com/avfs/avfs/vfs/osfs"
)

// ModeAppend is an alias of fs.ModeAppend.
//
// It is provided so callers can depend on go-service types while using standard
// file mode constants.
const ModeAppend = fs.ModeAppend

// DirEntry is an alias of fs.DirEntry.
//
// It represents a directory entry read from a directory.
type DirEntry = fs.DirEntry

// FileInfo is an alias of fs.FileInfo.
//
// It describes a file and is typically returned by Stat.
type FileInfo = fs.FileInfo

// FileMode is an alias of fs.FileMode.
//
// It represents a file's mode and permission bits.
type FileMode = fs.FileMode

// NewFS constructs an FS backed by the host OS filesystem.
//
// Internally this uses osfs.New() from github.com/avfs/avfs to provide an avfs.VFS
// implementation that delegates to the real operating system.
func NewFS() *FS {
	return &FS{VFS: osfs.New()}
}

// FS wraps an avfs.VFS and provides go-service-specific filesystem helpers.
//
// FS is intended to be used anywhere go-service needs filesystem access, while
// also providing:
//   - consistent path normalization (CleanPath / ExpandPath),
//   - convenience helpers (ExecutableName / ExecutableDir / PathExtension),
//   - and "source string" loading (ReadSource).
//
// The embedded avfs.VFS exposes a rich filesystem API; FS adds small, opinionated
// behavior on top (notably whitespace trimming for ReadFile/WriteFile).
type FS struct {
	avfs.VFS
}

// ReadFile reads the file at name after normalizing the path.
//
// The path is normalized using CleanPath, which expands a leading "~" (when
// present) and then cleans the path.
//
// The returned bytes are trimmed with bytes.TrimSpace, which is useful when
// reading configuration fragments or secrets where trailing newlines/whitespace
// are common.
func (fs *FS) ReadFile(name string) ([]byte, error) {
	b, err := fs.VFS.ReadFile(fs.CleanPath(name))

	return bytes.TrimSpace(b), err
}

// WriteFile writes data to name with perm after normalizing the path.
//
// The path is normalized using CleanPath, which expands a leading "~" (when
// present) and then cleans the path.
//
// The bytes written are trimmed with bytes.TrimSpace before being persisted. This
// mirrors ReadFile behavior and helps keep file-backed secrets/config fragments
// stable when inputs include trailing whitespace.
func (fs *FS) WriteFile(name string, data []byte, perm FileMode) error {
	return fs.VFS.WriteFile(fs.CleanPath(name), bytes.TrimSpace(data), perm)
}

// PathExists reports whether name exists in the filesystem.
//
// It returns false only when the underlying Stat returns an "not exist" error.
// Other errors (for example permission errors) are treated as "exists" by this
// function, because it cannot safely distinguish "missing" from "inaccessible"
// without inspecting the error further.
func (fs *FS) PathExists(name string) bool {
	if _, err := fs.Stat(name); fs.IsNotExist(err) {
		return false
	}

	return true
}

// IsNotExist reports whether err indicates a missing path.
//
// This helper uses errors.Is to match os.ErrNotExist.
func (*FS) IsNotExist(err error) bool {
	return errors.Is(err, os.ErrNotExist)
}

// PathExtension returns the file extension of path without the leading ".".
//
// If path has no extension, PathExtension returns an empty string.
//
// Examples:
//   - "file.yaml" -> "yaml"
//   - "file.test.yaml" -> "yaml"
//   - "file" -> ""
func (*FS) PathExtension(path string) string {
	e := filepath.Ext(path)
	if strings.IsEmpty(e) {
		return e
	}

	return e[1:]
}

// ExpandPath expands a leading "~" to the current user's home directory.
//
// If path is empty or does not start with "~", it is returned unchanged.
//
// The "~" expansion uses UserHomeDir, which panics if the home directory cannot
// be determined.
func (fs *FS) ExpandPath(path string) string {
	dir := UserHomeDir()
	if strings.IsAnyEmpty(dir, path) || path[0] != '~' {
		return path
	}

	return fs.Join(dir, path[1:])
}

// CleanPath expands "~" (when present) and then cleans the path.
//
// This is the normalization used by ReadFile and WriteFile to ensure consistent,
// user-friendly path handling across go-service.
func (fs *FS) CleanPath(path string) string {
	return fs.ExpandPath(fs.Clean(path))
}

// ExecutableName returns the base name of the running application executable.
//
// It uses Executable (which may panic if the executable path cannot be
// determined) and then returns the last path element.
func (fs *FS) ExecutableName() string {
	return fs.Base(Executable())
}

// ExecutableDir returns the directory containing the running application executable.
//
// It uses Executable (which may panic if the executable path cannot be
// determined) and then returns the directory portion of the path.
func (fs *FS) ExecutableDir() string {
	return fs.Dir(Executable())
}

// ReadSource reads a "source string" into bytes.
//
// ReadSource implements the go-service "source string" convention used to load
// configuration fragments and secrets from different sources without changing
// configuration structure.
//
// The supported forms are:
//
//   - "env:NAME" reads the value of environment variable NAME and returns it as bytes.
//   - "file:/path" reads the file at /path via ReadFile (including path cleaning and trimming).
//   - otherwise treats source as the literal value and returns it as bytes.
//
// Note: "env:" values are not trimmed; they are returned exactly as provided by
// the environment.
func (fs *FS) ReadSource(source string) ([]byte, error) {
	kind, path := strings.CutColon(source)
	switch kind {
	case "env":
		return strings.Bytes(Getenv(path)), nil
	case "file":
		return fs.ReadFile(path)
	default:
		return strings.Bytes(source), nil
	}
}
