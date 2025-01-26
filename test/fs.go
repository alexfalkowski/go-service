package test

import (
	"github.com/alexfalkowski/go-service/os"
)

// FS used for tests.
var FS = os.NewFS()

// BadFS for test.
type BadFS struct{}

func (f *BadFS) ReadFile(_ string) (string, error) {
	return "", ErrFailed
}

func (f *BadFS) PathExists(_ string) bool {
	return true
}

func (f *BadFS) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}
