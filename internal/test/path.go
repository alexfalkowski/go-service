package test

import "runtime"

// Path for test.
//
//nolint:dogsled
func Path(path string) string {
	_, b, _, _ := runtime.Caller(0)
	dir := FS.Dir(b)

	return FS.Join(dir, "../../test", path)
}

// FilePath for test.
func FilePath(path string) string {
	return "file:" + Path(path)
}
