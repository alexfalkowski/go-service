// Package types provides small, generic helper packages for working with common Go types.
//
// The packages under this namespace are intentionally lightweight and primarily exist to
// centralize small utilities that are reused across go-service without pulling in larger
// dependencies or rewriting the same helpers in multiple places.
//
// Subpackages:
//
//   - ptr: pointer helpers (for example constructing pointers to values and zero values).
//   - slices: slice helpers (for example appending conditionally and finding elements).
//   - structs: helpers for working with pointers and zero values (nil/empty/zero checks).
//
// # Design goals
//
// The helpers in these packages are designed to be:
//
//   - Generic: use Go generics to avoid reflection.
//   - Predictable: avoid surprising side effects; helpers generally return a derived value.
//   - Small: keep the surface area minimal and composable.
//
// This namespace is not intended to replace the standard library. Prefer the Go standard
// library first (for example `slices`, `maps`, `cmp`, `reflect` when necessary) and use
// these helpers when they improve readability or reduce repetitive boilerplate in
// go-service code.
package types
