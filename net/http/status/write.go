package status

import (
	"fmt"
	"net/http"
)

// WriteError will write the error to the response writer.
func WriteError(res http.ResponseWriter, err error) {
	writeError(res, err.Error(), Code(err))
}

func writeError(res http.ResponseWriter, message string, code int) {
	h := res.Header()
	h.Del("Content-Length")
	h.Set("Content-Type", "text/error; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")

	res.WriteHeader(code)
	fmt.Fprintln(res, message)
}
