package runtime

import (
	"fmt"
	"runtime"

	"github.com/alexfalkowski/go-service/v2/errors"
)

// MaxProcs sets the maximum number of CPUs that can execute Go code simultaneously and returns the previous setting.
//
// MaxProcs wraps the standard library's runtime.GOMAXPROCS so callers can stay on the project-owned runtime surface.
func MaxProcs(n int) int {
	return runtime.GOMAXPROCS(n)
}

// Caller reports file and line information for a stack frame.
//
// Caller wraps the standard library's runtime.Caller so callers can stay on the project-owned runtime surface while
// preserving its skip semantics.
func Caller(skip int) (pc uintptr, file string, line int, ok bool) {
	// Account for this wrapper's frame so skip still addresses the caller's stack.
	return runtime.Caller(skip + 1)
}

// ErrRecovered is a sentinel error used to mark errors produced by ConvertRecover.
//
// ConvertRecover wraps arbitrary panic values into an error and includes ErrRecovered
// in the returned error chain so callers can detect recovered-panics with:
//
//	errors.Is(err, runtime.ErrRecovered) == true
//
// When the recovered value is already an error, ConvertRecover wraps it so it remains
// accessible via [errors.As].
var ErrRecovered = errors.New("recovered")

// Must panics if err is non-nil.
//
// Must is intended for code paths where an error is not meaningfully recoverable,
// such as mandatory startup/configuration wiring. It is commonly used to reduce
// boilerplate when a function returns (T, error) and failure should abort:
//
//	v, err := build()
//	runtime.Must(err)
//
// Note: Must does not attach additional context. If you need context, wrap the
// error before calling Must.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// ConvertRecover converts a recovered panic value into an error wrapped with ErrRecovered.
//
// This helper is intended to be used with recover() inside a deferred function:
//
//	func run() (err error) {
//		defer func() {
//			if v := recover(); v != nil {
//				err = runtime.ConvertRecover(v)
//			}
//		}()
//		// ...
//		return nil
//	}
//
// The returned error always includes ErrRecovered in its chain. The recovered value is
// represented as:
//
//   - error: wrapped with %w (preserving the original error for [errors.As] / [errors.Is])
//   - string: included as text
//   - any other value: formatted with %v
func ConvertRecover(value any) error {
	switch recovered := value.(type) {
	case error:
		return fmt.Errorf("%w: %w", ErrRecovered, recovered)
	case string:
		return fmt.Errorf("%w: %s", ErrRecovered, recovered)
	default:
		return fmt.Errorf("%w: %v", ErrRecovered, recovered)
	}
}
