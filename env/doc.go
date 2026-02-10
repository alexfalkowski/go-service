// Package env provides service identity values derived from environment variables and defaults.
//
// This package defines small types and constructors for common service identity fields such as:
// service name, version, id, user id, and user agent.
//
// Constructors typically prefer an environment variable override (for example SERVICE_NAME, SERVICE_VERSION,
// SERVICE_ID, SERVICE_USER_ID) and otherwise fall back to a derived default.
package env
