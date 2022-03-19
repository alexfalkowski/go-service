package encoder

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// Request to be encoded.
func Request(req *http.Request) string {
	if req.Body == nil {
		return ""
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return ""
	}

	req.Body = io.NopCloser(bytes.NewBuffer(body))

	if !json.Valid(body) {
		return ""
	}

	return string(body)
}

// Response to be encoded.
func Response(resp *http.Response) string {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	if !json.Valid(body) {
		return ""
	}

	return string(body)
}
