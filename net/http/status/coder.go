package status

// Coder lets errors expose the HTTP status code that should be returned.
type Coder interface {
	// Code returns the HTTP status code to return, for example http.StatusNotFound.
	Code() int
}
