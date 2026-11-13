package errors_test

import (
	"testing"

	encoding "github.com/alexfalkowski/go-service/v3/encoding/errors"
	"github.com/alexfalkowski/go-service/v3/errors"
	"github.com/alexfalkowski/go-service/v3/io"
	"github.com/stretchr/testify/require"
)

func TestTrailingDataAcceptsEOFAndRejectsExtraOrInvalidData(t *testing.T) {
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
