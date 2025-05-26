package os

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/avfs/avfs"
	"github.com/avfs/avfs/vfs/osfs"
)

// ModeAppend is an alias to os.ModeAppend.
const ModeAppend = fs.ModeAppend

// FileMode is an alias to os.FileMode.
type FileMode = fs.FileMode

// NewFS for os.
func NewFS() *FS {
	return &FS{VFS: osfs.New()}
}

// FS for os.
type FS struct {
	avfs.VFS
}

// ReadFile for the name provided.
func (fs *FS) ReadFile(name string) ([]byte, error) {
	b, err := fs.VFS.ReadFile(fs.CleanPath(name))

	return bytes.TrimSpace(b), err
}

// WriteFile writes data to name with perm.
func (fs *FS) WriteFile(name string, data []byte, perm FileMode) error {
	return fs.VFS.WriteFile(fs.CleanPath(name), bytes.TrimSpace(data), perm)
}

// PathExists for os.
func (fs *FS) PathExists(name string) bool {
	if _, err := fs.Stat(name); fs.IsNotExist(err) {
		return false
	}

	return true
}

// IsNotExist for os.
func (*FS) IsNotExist(err error) bool {
	return errors.Is(err, os.ErrNotExist)
}

// PathExtension of the specified path.
func (*FS) PathExtension(path string) string {
	e := filepath.Ext(path)
	if strings.IsEmpty(e) {
		return e
	}

	return e[1:]
}

// ExpandPath will append the home dir if path starts with ~.
func (fs *FS) ExpandPath(path string) string {
	dir := UserHomeDir()
	if strings.IsEmpty(dir) || len(path) == 0 || path[0] != '~' {
		return path
	}

	return fs.Join(dir, path[1:])
}

// CleanPath makes sure that the path is expanded and in a clean format.
func (fs *FS) CleanPath(path string) string {
	return fs.ExpandPath(fs.Clean(path))
}

// ExecutableName of the running application.
func (fs *FS) ExecutableName() string {
	return fs.Base(Executable())
}

// ExecutableDir of the running application.
func (fs *FS) ExecutableDir() string {
	return fs.Dir(Executable())
}

// ReadSource will look at the source and read depending on env, file, etc.
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
