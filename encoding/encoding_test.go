package encoding_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestEncoder(t *testing.T) {
	for _, k := range test.Encoder.Keys() {
		require.NotNil(t, test.Encoder.Get(k))
	}

	for _, k := range []string{"test", "bob"} {
		require.Nil(t, test.Encoder.Get(k))
	}
}
