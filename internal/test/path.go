package test

import (
	"path/filepath"
	"runtime"
)

// Path for test.
//
//nolint:dogsled
func Path(path string) string {
	_, b, _, _ := runtime.Caller(0)
	dir := filepath.Dir(b)

	return filepath.Join(dir, "../../test", path)
}

// FilePath for test.
func FilePath(path string) string {
	return "file:" + Path(path)
}
