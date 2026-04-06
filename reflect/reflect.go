package reflect

import "reflect"

// IsNil reports whether value is nil.
//
// It returns true for both:
//   - a nil interface value, and
//   - an interface value holding a typed nil for kinds that can be nil
//     (chan, func, interface, map, pointer, and slice).
//
// For all other non-nil values, IsNil returns false.
func IsNil(value any) bool {
	if value == nil {
		return true
	}

	rv := reflect.ValueOf(value)

	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}
