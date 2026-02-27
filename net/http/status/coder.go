package status

// Coder allows errors to expose the HTTP status code that should be returned to a client.
//
// This interface is used by helpers in this package (for example Code and FromError) to extract a
// stable HTTP status code from an error value.
//
// Implementations should return a valid HTTP status code (e.g. 400, 404, 500). Callers typically use
// this in HTTP handlers/middleware to decide which status code to write for a given error.
type Coder interface {
	// Code returns the HTTP status code that should be returned, for example http.StatusNotFound.
	Code() int
}
