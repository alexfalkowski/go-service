package opentracing

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func encodeRequest(req *http.Request) string {
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

func encodeResponse(resp *http.Response) string {
	if resp.Body == nil {
		return ""
	}

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
