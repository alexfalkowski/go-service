package os

import (
	"encoding/base64"
	"os"
	"path/filepath"
)

// ReadBase64File for os.
func ReadBase64File(path string) ([]byte, error) {
	b, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(string(b))
}
