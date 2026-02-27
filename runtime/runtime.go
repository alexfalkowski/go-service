package runtime

import (
	"fmt"

	"github.com/alexfalkowski/go-service/v2/errors"
)

// ErrRecovered is a sentinel error used to mark errors produced by ConvertRecover.
//
// ConvertRecover wraps arbitrary panic values into an error and includes ErrRecovered
// in the returned error chain so callers can detect recovered-panics with:
//
//	errors.Is(err, runtime.ErrRecovered) == true
//
// When the recovered value is already an error, ConvertRecover wraps it so it remains
// accessible via errors.As.
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
//   - error: wrapped with %w (preserving the original error for errors.As / errors.Is)
//   - string: included as text
//   - any other value: formatted with %v
func ConvertRecover(value any) error {
	switch kind := value.(type) {
	case error:
		return fmt.Errorf("%w: %w", ErrRecovered, kind)
	case string:
		return fmt.Errorf("%w: %s", ErrRecovered, kind)
	default:
		return fmt.Errorf("%w: %v", ErrRecovered, kind)
	}
}
