package test

import "github.com/alexfalkowski/go-service/os"

// FS used for tests.
var FS = os.NewFS()

// ErrFS for test.
type ErrFS struct{}

func (f *ErrFS) ReadFile(_ string) ([]byte, error) {
	return nil, ErrFailed
}

func (*ErrFS) WriteFile(_ string, _ []byte, _ os.FileMode) error {
	return ErrFailed
}

func (f *ErrFS) PathExists(_ string) bool {
	return true
}

func (f *ErrFS) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}
