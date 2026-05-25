// Package psutil provides a debug endpoint that returns host/process system statistics.
//
// This package integrates github.com/shirou/gopsutil to collect a point-in-time snapshot of basic
// system metrics such as CPU, host info, load averages, memory/swap stats, and network I/O counters.
//
// The endpoint is typically registered on the debug HTTP mux under:
//
//	/debug/psutil
//
// (namespaced by service name via debug/http.Pattern).
//
// Notes:
//   - Collection is best-effort: errors returned by gopsutil calls are intentionally ignored and the
//     corresponding fields may be partially populated or empty.
//   - The endpoint is intended for diagnostics and operations; consider access control and exposure
//     carefully in production environments.
package psutil
