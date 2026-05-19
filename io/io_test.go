package io_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/stretchr/testify/require"
)

func TestWriteString(t *testing.T) {
	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	n, err := io.WriteString(buffer, "hello")
	require.NoError(t, err)
	require.Equal(t, 5, n)
	require.Equal(t, "hello", buffer.String())
}
