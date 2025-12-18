package time_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestValidNTP(t *testing.T) {
	c := &time.Config{Kind: "ntp", Address: "0.beevik-ntp.pool.ntp.org"}

	n, err := time.NewNetwork(c)
	require.NoError(t, err)

	_, err = n.Now()
	require.NoError(t, err)
}

func TestInvalidNTP(t *testing.T) {
	c := &time.Config{Kind: "ntp"}

	n, err := time.NewNetwork(c)
	require.NoError(t, err)

	_, err = n.Now()
	require.Error(t, err)
}
