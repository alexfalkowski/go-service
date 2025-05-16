package test

import (
	"github.com/alexfalkowski/go-service/os"
	"github.com/avfs/avfs"
	"github.com/avfs/avfs/vfs/failfs"
	"github.com/avfs/avfs/vfs/memfs"
)

func init() {
	f := failfs.New(memfs.New())
	_ = f.SetFailFunc(fail)

	ErrFS = &os.FS{VFS: f}
}

var (
	// FS used for tests.
	FS = os.NewFS()

	// ErrFS for tests.
	ErrFS *os.FS
)

func fail(_ avfs.VFSBase, _ avfs.FnVFS, _ *failfs.FailParam) error {
	return ErrFailed
}
