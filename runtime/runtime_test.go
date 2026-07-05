package runtime_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/stretchr/testify/require"
)

func TestMustPanicsWithError(t *testing.T) {
	require.Panics(t, func() { runtime.Must(test.ErrFailed) })
}

func TestRecover(t *testing.T) {
	tests := []struct {
		value   any
		name    string
		message string
	}{
		{name: "error", value: test.ErrFailed, message: "recovered: failed"},
		{name: "string", value: "test", message: "recovered: test"},
		{name: "int", value: 1, message: "recovered: 1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			func() {
				defer func() {
					if recovered := recover(); recovered != nil {
						err = runtime.ConvertRecover(recovered)
					}
				}()

				panic(tt.value)
			}()

			require.Error(t, err)
			require.Equal(t, tt.message, err.Error())
		})
	}
}
