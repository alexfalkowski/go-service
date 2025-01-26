package test

import (
	"github.com/alexfalkowski/go-service/os"
)

// FS used for tests.
var FS = os.NewFS()

// ErrFS for test.
type ErrFS struct{}

func (f *ErrFS) ReadFile(_ string) (string, error) {
	return "", ErrFailed
}

func (f *ErrFS) PathExists(_ string) bool {
	return true
}

func (f *ErrFS) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}
