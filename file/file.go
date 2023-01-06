package file

import (
	"path/filepath"
)

// Extension of file.
func Extension(f string) string {
	e := filepath.Ext(f)
	if e == "" {
		return e
	}

	return e[1:]
}
