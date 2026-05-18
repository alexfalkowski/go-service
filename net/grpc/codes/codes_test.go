package codes_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/stretchr/testify/require"
)

func TestStatusText(t *testing.T) {
	require.Equal(t, codes.Unauthenticated.String(), codes.StatusText(codes.Unauthenticated))
}
