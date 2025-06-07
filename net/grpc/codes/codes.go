package codes

import "google.golang.org/grpc/codes"

// Code is an alias for codes.Code.
type Code = codes.Code

var (
	// Aborted is an alias for codes.Aborted.
	Aborted = codes.Aborted

	// AlreadyExists is an alias for codes.AlreadyExists.
	AlreadyExists = codes.AlreadyExists

	// Canceled is an alias for codes.Canceled.
	Canceled = codes.Canceled

	// DataLoss is an alias for codes.DataLoss.
	DataLoss = codes.DataLoss

	// DeadlineExceeded is an alias for codes.DeadlineExceeded.
	DeadlineExceeded = codes.DeadlineExceeded

	// FailedPrecondition is an alias for codes.FailedPrecondition.
	FailedPrecondition = codes.FailedPrecondition

	// InvalidArgument is an alias for codes.InvalidArgument.
	InvalidArgument = codes.InvalidArgument

	// Internal is an alias for codes.Internal.
	Internal = codes.Internal

	// PermissionDenied is an alias for codes.PermissionDenied.
	PermissionDenied = codes.PermissionDenied

	// OK is an alias for codes.OK.
	OK = codes.OK

	// OutOfRange is an alias for codes.OutOfRange.
	OutOfRange = codes.OutOfRange

	// NotFound is an alias for codes.NotFound.
	NotFound = codes.NotFound

	// ResourceExhausted is an alias for codes.ResourceExhausted.
	ResourceExhausted = codes.ResourceExhausted

	// Unavailable is an alias for codes.Unavailable.
	Unavailable = codes.Unavailable

	// Unimplemented is an alias for codes.Unimplemented.
	Unimplemented = codes.Unimplemented

	// Unknown is an alias for codes.Unknown.
	Unknown = codes.Unknown

	// Unauthenticated is an alias for codes.Unauthenticated.
	Unauthenticated = codes.Unauthenticated
)
