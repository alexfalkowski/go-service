package gob_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/encoding/gob"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestEncoder(t *testing.T) {
	encoder := gob.NewEncoder()

	bytes := test.Pool.Get()
	defer test.Pool.Put(bytes)

	require.NoError(t, encoder.Encode(bytes, map[string]string{"test": "test"}))

	var msg map[string]string
	require.NoError(t, encoder.Decode(bytes, &msg))
	require.Equal(t, map[string]string{"test": "test"}, msg)
}
