package test

import (
	"path/filepath"
	"runtime"
	"strings"
)

// Path for test.
//
//nolint:dogsled
func Path(path string) string {
	_, b, _, _ := runtime.Caller(0)
	dir := filepath.Dir(b)

	if strings.HasPrefix(path, "certs") {
		return filepath.Join(dir, "../../test", path)
	}

	return filepath.Join(dir, path)
}

// FilePath for test.
func FilePath(path string) string {
	return "file:" + Path(path)
}
