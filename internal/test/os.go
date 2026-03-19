package test

import (
	"github.com/alexfalkowski/go-service/v2/os"
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
	// FS is the shared filesystem used by test helpers.
	FS = os.NewFS()

	// ErrFS is a filesystem test double whose operations fail with ErrFailed.
	ErrFS *os.FS
)

func fail(_ avfs.VFSBase, _ avfs.FnVFS, _ *failfs.FailParam) error {
	return ErrFailed
}
