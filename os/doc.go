// Package os provides filesystem and OS-related helpers used by go-service.
//
// This package contains a filesystem abstraction (`FS`) that supports the go-service "source string" pattern
// (see FS.ReadSource), along with small wrappers around standard library os functionality.
//
// Some helpers in this package are intentionally strict and will panic on unexpected OS errors (for example
// Executable, UserHomeDir, and UserConfigDir) via runtime.Must.
//
// Start with `FS` and its helpers.
package os
