package time_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestNil(t *testing.T) {
	net, err := time.NewNetwork(nil)
	require.NoError(t, err)
	require.Nil(t, net)
}

func TestInvalid(t *testing.T) {
	_, err := time.NewNetwork(&time.Config{Kind: "invalid"})
	require.Error(t, err)
}
