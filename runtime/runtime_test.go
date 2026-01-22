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

	functions := []fun{
		func() (err error) {
			defer func() {
				if r := recover(); r != nil {
					err = runtime.ConvertRecover(r)
					require.Equal(t, "recovered: failed", err.Error())
				}
			}()

			panic(test.ErrFailed)
		},
		func() (err error) {
			defer func() {
				if r := recover(); r != nil {
					err = runtime.ConvertRecover(r)
					require.Equal(t, "recovered: test", err.Error())
				}
			}()

			panic("test")
		},
		func() (err error) {
			defer func() {
				if r := recover(); r != nil {
					err = runtime.ConvertRecover(r)
					require.Equal(t, "recovered: 1", err.Error())
				}
			}()

			panic(1)
		},
	}

	for _, f := range functions {
		require.Error(t, f())
	}
}
