package media_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/stretchr/testify/require"
)

func TestParseInvalidMediaType(t *testing.T) {
	_, _, err := media.Parse("json")

	require.ErrorIs(t, err, media.ErrInvalidType)
	require.True(t, errors.Is(err, media.ErrInvalidType))
}
