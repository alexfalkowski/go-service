package codes

import "google.golang.org/grpc/codes"

// Code is the canonical gRPC status code type used by go-service.
//
// It is an alias of google.golang.org/grpc/codes.Code, re-exported so callers
// can depend on go-service packages while still using the upstream gRPC status
// code representation.
//
// Code values are most commonly produced/consumed via the sibling package
// net/grpc/status, for example:
//
//	err := status.Error(codes.NotFound, "widget does not exist")
//	c := status.Code(err) // c is a codes.Code
//
// For the authoritative meaning of each code and guidance on when to use them,
// consult the upstream gRPC documentation.
type Code = codes.Code

const (
	// Aborted indicates the operation was aborted, typically due to a concurrency
	// issue such as sequencer check failures, transaction aborts, etc.
	//
	// This is an alias for google.golang.org/grpc/codes.Aborted.
	Aborted = codes.Aborted

	// AlreadyExists indicates an attempt was made to create a resource that
	// already exists.
	//
	// This is an alias for google.golang.org/grpc/codes.AlreadyExists.
	AlreadyExists = codes.AlreadyExists

	// Canceled indicates the operation was canceled, typically by the caller.
	//
	// This is an alias for google.golang.org/grpc/codes.Canceled.
	Canceled = codes.Canceled

	// DataLoss indicates unrecoverable data loss or corruption.
	//
	// This is an alias for google.golang.org/grpc/codes.DataLoss.
	DataLoss = codes.DataLoss

	// DeadlineExceeded indicates the operation timed out before completion.
	//
	// This is an alias for google.golang.org/grpc/codes.DeadlineExceeded.
	DeadlineExceeded = codes.DeadlineExceeded

	// FailedPrecondition indicates the operation was rejected because the system
	// is not in a state required for the operation's execution.
	//
	// This is an alias for google.golang.org/grpc/codes.FailedPrecondition.
	FailedPrecondition = codes.FailedPrecondition

	// InvalidArgument indicates the client specified an invalid argument. Note
	// that this differs from FailedPrecondition: InvalidArgument indicates
	// arguments that are problematic regardless of the state of the system.
	//
	// This is an alias for google.golang.org/grpc/codes.InvalidArgument.
	InvalidArgument = codes.InvalidArgument

	// Internal indicates an internal error. This generally means some invariants
	// expected by the underlying system have been broken.
	//
	// This is an alias for google.golang.org/grpc/codes.Internal.
	Internal = codes.Internal

	// PermissionDenied indicates the caller does not have permission to execute
	// the specified operation.
	//
	// This is an alias for google.golang.org/grpc/codes.PermissionDenied.
	PermissionDenied = codes.PermissionDenied

	// OK indicates the operation completed successfully.
	//
	// This is an alias for google.golang.org/grpc/codes.OK.
	OK = codes.OK

	// OutOfRange indicates an operation was attempted past the valid range. This
	// typically refers to range errors on input values.
	//
	// This is an alias for google.golang.org/grpc/codes.OutOfRange.
	OutOfRange = codes.OutOfRange

	// NotFound indicates the requested entity was not found.
	//
	// This is an alias for google.golang.org/grpc/codes.NotFound.
	NotFound = codes.NotFound

	// ResourceExhausted indicates some resource has been exhausted, perhaps a
	// per-user quota, or perhaps the entire file system is out of space.
	//
	// This is an alias for google.golang.org/grpc/codes.ResourceExhausted.
	ResourceExhausted = codes.ResourceExhausted

	// Unavailable indicates the service is currently unavailable. This is most
	// often a transient condition and may be corrected by retrying with backoff.
	//
	// This is an alias for google.golang.org/grpc/codes.Unavailable.
	Unavailable = codes.Unavailable

	// Unimplemented indicates the operation is not implemented or is not
	// supported/enabled in this service.
	//
	// This is an alias for google.golang.org/grpc/codes.Unimplemented.
	Unimplemented = codes.Unimplemented

	// Unknown indicates an unknown error. This is generally returned when an
	// error does not map cleanly to any other gRPC code.
	//
	// This is an alias for google.golang.org/grpc/codes.Unknown.
	Unknown = codes.Unknown

	// Unauthenticated indicates the request does not have valid authentication
	// credentials for the operation.
	//
	// This is an alias for google.golang.org/grpc/codes.Unauthenticated.
	Unauthenticated = codes.Unauthenticated
)
