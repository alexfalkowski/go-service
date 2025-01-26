package os

// FileSystem borrows concepts from io/fs.
type FileSystem interface {
	// ReadFile for the path provided.
	ReadFile(path string) (string, error)

	// PathExists for the path provided.
	PathExists(path string) bool

	// IsNotExist whether the error is os.ErrNotExist.
	IsNotExist(err error) bool
}

// NewFS for os.
func NewFS() FileSystem {
	return &fs{}
}

type fs struct{}

func (f *fs) ReadFile(path string) (string, error) {
	return ReadFile(path)
}

func (f *fs) PathExists(name string) bool {
	return PathExists(name)
}

func (f *fs) IsNotExist(err error) bool {
	return IsNotExist(err)
}
