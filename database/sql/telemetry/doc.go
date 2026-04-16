// Package telemetry exposes selected otelsql helpers through the go-service
// SQL import tree.
//
// This package keeps repository SQL instrumentation on a go-service import path
// while preserving otelsql behavior for opening instrumented connections,
// wrapping drivers, and registering DB stats metrics.
package telemetry
