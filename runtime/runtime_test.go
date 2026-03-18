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
	functions := []struct {
		fun  func(t *testing.T) (err error)
		name string
	}{
		{
			name: "error",
			fun: func(t *testing.T) (err error) {
				t.Helper()

				defer func() {
					if r := recover(); r != nil {
						err = runtime.ConvertRecover(r)
						require.Equal(t, "recovered: failed", err.Error())
					}
				}()

				panic(test.ErrFailed)
			},
		},
		{
			name: "string",
			fun: func(t *testing.T) (err error) {
				t.Helper()

				defer func() {
					if r := recover(); r != nil {
						err = runtime.ConvertRecover(r)
						require.Equal(t, "recovered: test", err.Error())
					}
				}()

				panic("test")
			},
		},
		{
			name: "int",
			fun: func(t *testing.T) (err error) {
				t.Helper()

				defer func() {
					if r := recover(); r != nil {
						err = runtime.ConvertRecover(r)
						require.Equal(t, "recovered: 1", err.Error())
					}
				}()

				panic(1)
			},
		},
	}

	for _, tt := range functions {
		t.Run(tt.name, func(t *testing.T) {
			require.Error(t, tt.fun(t))
		})
	}
}
