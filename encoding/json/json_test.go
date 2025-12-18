package json_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding/json"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestEncode(t *testing.T) {
	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	encoder := json.NewEncoder()
	msg := map[string]string{"test": "test"}

	require.NoError(t, encoder.Encode(buffer, msg))
	require.JSONEq(t, "{\"test\":\"test\"}", bytes.String(bytes.TrimSpace(test.Pool.Copy(buffer))))
}

func TestDecode(t *testing.T) {
	encoder := json.NewEncoder()
	var msg map[string]string

	require.NoError(t, encoder.Decode(bytes.NewBufferString("{\"test\":\"test\"}"), &msg))
	require.Equal(t, map[string]string{"test": "test"}, msg)
}
