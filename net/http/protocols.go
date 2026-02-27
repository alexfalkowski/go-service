package http

import "net/http"

// Protocols constructs an http.Protocols value enabling the HTTP protocols supported by go-service.
//
// The returned configuration enables:
//   - HTTP/1.1
//   - HTTP/2 (when negotiated, typically over TLS)
//   - h2c (unencrypted HTTP/2)
//
// This helper is used by go-service HTTP server and transport construction to consistently enable
// the same protocol set across servers and clients.
func Protocols() *http.Protocols {
	protocols := &http.Protocols{}
	protocols.SetHTTP1(true)
	protocols.SetHTTP2(true)
	protocols.SetUnencryptedHTTP2(true)

	return protocols
}
