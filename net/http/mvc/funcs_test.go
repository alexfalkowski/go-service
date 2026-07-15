package mvc_test

import (
	"log/slog"
	"testing"

	"github.com/alexfalkowski/go-service/v2/net/http/mvc"
	"github.com/stretchr/testify/require"
)

func TestNewFunctionMapRemovesUnsafeShuffle(t *testing.T) {
	fmap := mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()})

	require.NotContains(t, fmap, "shuffle")
	require.NotContains(t, fmap, "safeShuffle")
}
