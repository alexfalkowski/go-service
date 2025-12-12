package runtime_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/stretchr/testify/require"
)

func TestPanic(t *testing.T) {
	require.Panics(t, func() { runtime.Must(test.ErrFailed) })
}

func TestRecover(t *testing.T) {
	type fun func() (err error)

	errPanic := func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = runtime.ConvertRecover(r)
			}
		}()

		panic(test.ErrFailed)
	}

	strPanic := func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = runtime.ConvertRecover(r)
			}
		}()

		panic("test")
	}

	intPanic := func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = runtime.ConvertRecover(r)
			}
		}()

		panic(1)
	}

	for _, f := range []fun{errPanic, strPanic, intPanic} {
		require.Error(t, f())
	}
}
