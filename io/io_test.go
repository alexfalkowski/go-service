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

	value := "hello"
	written, err := io.WriteString(buffer, value)
	require.NoError(t, err)
	require.Equal(t, len(value), written)
	require.Equal(t, value, buffer.String())
}
