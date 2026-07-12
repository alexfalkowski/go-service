// Package method defines gRPC method policy helpers.
//
// Policy classifies methods before the server begins serving. Under the
// standard go-service gRPC server stack, the classifications have these
// middleware consequences:
//
//   - Normal methods participate in metadata extraction, logging, and any
//     configured token verification, limiting, and access control.
//   - Unauthenticated methods retain metadata extraction, logging, and
//     limiting, but bypass token verification and access control.
//   - Operation unary methods bypass metadata extraction, logging, token
//     verification, limiting, and access control.
//   - Operation streams retain metadata extraction and stream limiting, but
//     bypass logging, token verification, and access control.
//
// Operation streams intentionally remain limited because long-lived streams,
// such as gRPC health Watch, can hold server resources. Populate a [Policy]
// during service registration; it is not a runtime reconfiguration API.
package method
