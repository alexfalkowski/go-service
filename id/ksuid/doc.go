// Package ksuid provides KSUID generation and parsing helpers used by go-service.
//
// This package integrates KSUID generation behind the go-service ID abstraction.
// KSUID values are operational identifiers with time-sortable characteristics;
// they are not secret material or bearer tokens.
//
// Use [NewGenerator] to generate KSUID strings and [Parse] to decode them into a [KSUID].
package ksuid
