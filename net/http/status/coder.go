package status

// Coder allows errors to implement so we can return the code needed.
type Coder interface {
	// Code reflects the status code to return, e.g: http.StatusNotFound.
	Code() int
}
