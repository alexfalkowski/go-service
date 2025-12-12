package time_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestMustParseDuration(t *testing.T) {
	require.Panics(t, func() { time.MustParseDuration("test") })
}
