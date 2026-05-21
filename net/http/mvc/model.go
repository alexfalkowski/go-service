package mvc

// Error is the client-safe model passed to MVC error views.
type Error struct {
	// Message is the safe, client-visible error message.
	Message string

	// Code is the HTTP status code associated with the error.
	Code int
}
