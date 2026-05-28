package hjson_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding/hjson"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
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

func TestEncodeReturnsError(t *testing.T) {
	encoder := hjson.NewEncoder()
	buf := &bytes.Buffer{}

	require.Error(t, encoder.Encode(buf, func() {}))
}

func TestEncodeReturnsWriteError(t *testing.T) {
	encoder := hjson.NewEncoder()

	require.ErrorIs(t, encoder.Encode(test.ErrWriter{}, &message{Test: "test"}), test.ErrFailed)
}

func TestDecode(t *testing.T) {
	encoder := hjson.NewEncoder()
	msg := &message{}

	require.NoError(t, encoder.Decode(bytes.NewBufferString("{\n  // hjson comment\n  test: test\n}\n"), msg))
	require.Equal(t, "test", msg.Test)
}

func TestDecodeReturnsReadError(t *testing.T) {
	encoder := hjson.NewEncoder()
	msg := &message{}

	require.ErrorIs(t, encoder.Decode(&test.ErrReaderCloser{}, msg), test.ErrFailed)
}

func TestDecodeRejectsDuplicateKeys(t *testing.T) {
	encoder := hjson.NewEncoder()
	msg := &message{}

	err := encoder.Decode(bytes.NewBufferString("{\n  test: first\n  test: second\n}\n"), msg)
	require.Error(t, err)
	require.Contains(t, strings.ToLower(err.Error()), "duplicate")
}
