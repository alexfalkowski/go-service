package http

import "net/http"

// Protocols constructs an [http.Protocols] value enabling the HTTP protocols supported by go-service.
//
// The returned configuration enables:
//   - HTTP/1.1
//   - HTTP/2 (when negotiated, typically over TLS)
//   - h2c (unencrypted HTTP/2)
//
// This helper is used by go-service HTTP server and transport construction to consistently enable
// the same protocol set across servers and clients.
//
// Client-side h2c caveat: enabling both HTTP/1.1 and h2c on a client transport, as this helper does,
// does not make [net/http.Transport] negotiate h2c against a plain "http://" URL — with HTTP/1.1 also
// enabled, the standard library client uses HTTP/1.1 for cleartext requests. h2c prior-knowledge
// negotiation only happens when a transport has HTTP/1.1 disabled (`SetHTTP1(false)`) so h2c is the only
// enabled cleartext option. This helper's default (both enabled) intentionally favors broad
// interoperability with HTTP/1.1-only servers over automatic h2c negotiation to arbitrary endpoints.
func Protocols() *http.Protocols {
	protocols := &http.Protocols{}
	protocols.SetHTTP1(true)
	protocols.SetHTTP2(true)
	protocols.SetUnencryptedHTTP2(true)

	return protocols
}
