package codes

import "google.golang.org/grpc/codes"

// OK is an alias for [google.golang.org/grpc/codes.OK].
const OK = codes.OK

// Canceled is an alias for [google.golang.org/grpc/codes.Canceled].
const Canceled = codes.Canceled

// Unknown is an alias for [google.golang.org/grpc/codes.Unknown].
const Unknown = codes.Unknown

// InvalidArgument is an alias for [google.golang.org/grpc/codes.InvalidArgument].
const InvalidArgument = codes.InvalidArgument

// DeadlineExceeded is an alias for [google.golang.org/grpc/codes.DeadlineExceeded].
const DeadlineExceeded = codes.DeadlineExceeded

// NotFound is an alias for [google.golang.org/grpc/codes.NotFound].
const NotFound = codes.NotFound

// AlreadyExists is an alias for [google.golang.org/grpc/codes.AlreadyExists].
const AlreadyExists = codes.AlreadyExists

// PermissionDenied is an alias for [google.golang.org/grpc/codes.PermissionDenied].
const PermissionDenied = codes.PermissionDenied

// ResourceExhausted is an alias for [google.golang.org/grpc/codes.ResourceExhausted].
const ResourceExhausted = codes.ResourceExhausted

// FailedPrecondition is an alias for [google.golang.org/grpc/codes.FailedPrecondition].
const FailedPrecondition = codes.FailedPrecondition

// Aborted is an alias for [google.golang.org/grpc/codes.Aborted].
const Aborted = codes.Aborted

// OutOfRange is an alias for [google.golang.org/grpc/codes.OutOfRange].
const OutOfRange = codes.OutOfRange

// Unimplemented is an alias for [google.golang.org/grpc/codes.Unimplemented].
const Unimplemented = codes.Unimplemented

// Internal is an alias for [google.golang.org/grpc/codes.Internal].
const Internal = codes.Internal

// Unavailable is an alias for [google.golang.org/grpc/codes.Unavailable].
const Unavailable = codes.Unavailable

// DataLoss is an alias for [google.golang.org/grpc/codes.DataLoss].
const DataLoss = codes.DataLoss

// Unauthenticated is an alias for [google.golang.org/grpc/codes.Unauthenticated].
const Unauthenticated = codes.Unauthenticated

// Code is the canonical gRPC status code type used by go-service.
//
// It is an alias of [google.golang.org/grpc/codes.Code], re-exported so callers
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

// StatusText returns the standard gRPC status text for the code.
func StatusText(code Code) string {
	return code.String()
}
