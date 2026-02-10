// Package runtime provides runtime-related helpers used by go-service.
//
// This package contains small utilities such as:
//   - Must, which panics on non-nil errors (used for strict startup/config paths),
//   - ConvertRecover, which converts recovered panic values into errors, and
//   - Version, which reports build info version (or "development").
//
// It also integrates optional runtime configuration such as setting Go's memory limit (see RegisterMemLimit).
//
// Start with `Must`, `ConvertRecover`, `Version`, and `Module`.
package runtime
