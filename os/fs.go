package os

// FileSystem borrows concepts from io/fs.
type FileSystem interface {
	// ReadFile for the path provided.
	ReadFile(path string) (string, error)

	// FileExists for the path provided.
	FileExists(path string) bool
}

// NewFS for os.
func NewFS() FileSystem {
	return &fs{}
}

type fs struct{}

func (f *fs) ReadFile(path string) (string, error) {
	return ReadFile(path)
}

func (f *fs) FileExists(name string) bool {
	return FileExists(name)
}
