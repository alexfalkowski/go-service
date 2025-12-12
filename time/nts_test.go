package time_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestValidNTS(t *testing.T) {
	c := &time.Config{Kind: "nts", Address: "time.cloudflare.com"}

	n, err := time.NewNetwork(c)
	require.NoError(t, err)

	_, err = n.Now()
	require.NoError(t, err)
}

func TestInvalidNTS(t *testing.T) {
	c := &time.Config{Kind: "nts"}

	n, err := time.NewNetwork(c)
	require.NoError(t, err)

	_, err = n.Now()
	require.Error(t, err)
}
