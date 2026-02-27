// Package env provides service identity values derived from environment variables and defaults.
//
// This package defines small types and constructors for common identity fields that are used across
// go-service for consistent naming/versioning and outbound metadata, such as:
//   - service name (Name)
//   - service version (Version)
//   - service instance id (ID)
//   - service user id (UserID)
//   - HTTP User-Agent value (UserAgent)
//
// # Environment variable overrides
//
// Constructors in this package typically prefer an environment variable override and otherwise fall
// back to a derived default:
//
//   - SERVICE_NAME: overrides the service name (otherwise executable name)
//   - SERVICE_VERSION: overrides the service version (otherwise build/runtime version)
//   - SERVICE_ID: overrides the service instance id (otherwise generated)
//   - SERVICE_USER_ID: overrides the service user id (otherwise service name)
//
// # Conventions
//
// Identity values are represented as small string wrapper types with a `String()` method.
// These wrappers preserve the underlying semantics while keeping imports consistent across go-service.
//
// Start with `NewName`, `NewVersion`, `NewID`, `NewUserID`, and `NewUserAgent`.
package env
