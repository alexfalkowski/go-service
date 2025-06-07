package status

import "google.golang.org/grpc/status"

var (
	// Code is an alias for status.Code.
	Code = status.Code

	// Error is an alias for status.Error.
	Error = status.Error

	// Errorf is an alias for status.Errorf.
	Errorf = status.Errorf
)
