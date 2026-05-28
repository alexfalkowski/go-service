package errors_test

import (
	"testing"

	encoding "github.com/alexfalkowski/go-service/v2/encoding/errors"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/stretchr/testify/require"
)

func TestTrailingData(t *testing.T) {
	parseErr := errors.New("parse failed")

	t.Run("eof", func(t *testing.T) {
		require.NoError(t, encoding.TrailingData(io.EOF))
	})

	t.Run("extra value", func(t *testing.T) {
		require.ErrorIs(t, encoding.TrailingData(nil), encoding.ErrTrailingData)
	})

	t.Run("parse error", func(t *testing.T) {
		err := encoding.TrailingData(parseErr)
		require.ErrorIs(t, err, encoding.ErrTrailingData)
		require.ErrorIs(t, err, parseErr)
	})
}
