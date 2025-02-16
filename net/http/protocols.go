package http

import (
	"net/http"
)

// Protocols returns all the protocols supported.
func Protocols() *http.Protocols {
	protocols := &http.Protocols{}
	protocols.SetHTTP1(true)
	protocols.SetHTTP2(true)
	protocols.SetUnencryptedHTTP2(true)

	return protocols
}
