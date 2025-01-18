package io

import (
	"bytes"
	"io"
)

// ReadAll reads all the bytes from the io.Reader and returns the bytes with an io.ReadCloser.
func ReadAll(r io.Reader) ([]byte, io.ReadCloser, error) {
	payload, err := io.ReadAll(r)
	if err != nil {
		return nil, nil, err
	}

	return payload, io.NopCloser(bytes.NewReader(payload)), nil
}
