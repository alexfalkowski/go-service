package hjson_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding/hjson"
	"github.com/stretchr/testify/require"
)

type message struct {
	Test string `json:"test"`
}

func TestEncode(t *testing.T) {
	encoder := hjson.NewEncoder()
	buf := &bytes.Buffer{}

	require.NoError(t, encoder.Encode(buf, &message{Test: "test"}))
	require.Contains(t, buf.String(), "test")
}

func TestDecode(t *testing.T) {
	encoder := hjson.NewEncoder()
	msg := &message{}

	require.NoError(t, encoder.Decode(bytes.NewBufferString("{\n  // hjson comment\n  test: test\n}\n"), msg))
	require.Equal(t, "test", msg.Test)
}
