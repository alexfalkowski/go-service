package test

import "runtime"

// Path resolves a fixture path relative to the repository's `test/` directory.
//
//nolint:dogsled
func Path(path string) string {
	_, b, _, _ := runtime.Caller(0)
	dir := FS.Dir(b)

	return FS.Join(dir, "../../test", path)
}

// FilePath prefixes Path with `file:` so it can be consumed by source-string config fields.
func FilePath(path string) string {
	return "file:" + Path(path)
}
