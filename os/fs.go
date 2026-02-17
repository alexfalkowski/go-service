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

// ModeAppend is an alias to os.ModeAppend.
const ModeAppend = fs.ModeAppend

// DirEntry is an alias for fs.DirEntry.
type DirEntry = fs.DirEntry

// FileInfo is an alias for fs.FileInfo.
type FileInfo = fs.FileInfo

// FileMode is an alias for fs.FileMode.
type FileMode = fs.FileMode

// NewFS constructs a new filesystem wrapper backed by the OS filesystem.
func NewFS() *FS {
	return &FS{VFS: osfs.New()}
}

// FS wraps an avfs.VFS and provides go-service-specific filesystem helpers.
type FS struct {
	avfs.VFS
}

// ReadFile reads the file at name after cleaning the path.
//
// The returned bytes are trimmed with bytes.TrimSpace.
func (fs *FS) ReadFile(name string) ([]byte, error) {
	b, err := fs.VFS.ReadFile(fs.CleanPath(name))

	return bytes.TrimSpace(b), err
}

// WriteFile writes data to name with perm after cleaning the path.
//
// The written bytes are trimmed with bytes.TrimSpace.
func (fs *FS) WriteFile(name string, data []byte, perm FileMode) error {
	return fs.VFS.WriteFile(fs.CleanPath(name), bytes.TrimSpace(data), perm)
}

// PathExists reports whether name exists in the filesystem.
func (fs *FS) PathExists(name string) bool {
	if _, err := fs.Stat(name); fs.IsNotExist(err) {
		return false
	}

	return true
}

// IsNotExist reports whether err indicates a missing path.
func (*FS) IsNotExist(err error) bool {
	return errors.Is(err, os.ErrNotExist)
}

// PathExtension returns the file extension of path without the leading ".".
//
// If path has no extension, it returns an empty string.
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
func (fs *FS) ExpandPath(path string) string {
	dir := UserHomeDir()
	if strings.IsAnyEmpty(dir, path) || path[0] != '~' {
		return path
	}

	return fs.Join(dir, path[1:])
}

// CleanPath expands "~" (when present) and then cleans the path.
func (fs *FS) CleanPath(path string) string {
	return fs.ExpandPath(fs.Clean(path))
}

// ExecutableName returns the base name of the running application executable.
func (fs *FS) ExecutableName() string {
	return fs.Base(Executable())
}

// ExecutableDir returns the directory of the running application executable.
func (fs *FS) ExecutableDir() string {
	return fs.Dir(Executable())
}

// ReadSource reads a "source string" into bytes.
//
// The supported forms are:
//   - "env:NAME"  -> reads the value of environment variable NAME
//   - "file:/path" -> reads the file at /path (via ReadFile, including path cleaning and trimming)
//   - otherwise   -> treats source as the literal value
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
