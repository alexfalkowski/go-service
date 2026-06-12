package json_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding/errors"
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

func TestEncodeReturnsError(t *testing.T) {
	buffer := test.Pool.Get()
	defer test.Pool.Put(buffer)

	encoder := json.NewEncoder()

	require.Error(t, encoder.Encode(buffer, func() {}))
}

func TestDecode(t *testing.T) {
	encoder := json.NewEncoder()
	var msg map[string]string

	require.NoError(t, encoder.Decode(bytes.NewBufferString("{\"test\":\"test\"}"), &msg))
	require.Equal(t, map[string]string{"test": "test"}, msg)
}

func TestDecodeAcceptsTrailingWhitespace(t *testing.T) {
	encoder := json.NewEncoder()
	var msg map[string]string

	require.NoError(t, encoder.Decode(bytes.NewBufferString("{\"test\":\"test\"} \n\t"), &msg))
	require.Equal(t, map[string]string{"test": "test"}, msg)
}

func TestDecodeRejectsUnknownFields(t *testing.T) {
	encoder := json.NewEncoder()
	msg := &message{}

	err := encoder.Decode(bytes.NewBufferString("{\"test\":\"test\",\"extra\":\"ignored\"}"), msg)

	require.Error(t, err)
	require.Contains(t, err.Error(), "extra")
}

func TestDecodeRejectsTrailingData(t *testing.T) {
	for _, input := range []string{
		`{"test":"test"} garbage`,
		`{"test":"test"}{"test":"other"}`,
	} {
		t.Run(input, func(t *testing.T) {
			encoder := json.NewEncoder()
			var msg map[string]string

			err := encoder.Decode(bytes.NewBufferString(input), &msg)

			require.ErrorIs(t, err, errors.ErrTrailingData)
		})
	}
}

func TestMarshal(t *testing.T) {
	msg := map[string]string{"test": "test"}

	data, err := json.Marshal(msg)
	require.NoError(t, err)
	require.JSONEq(t, "{\"test\":\"test\"}", string(data))
	require.Contains(t, string(data), "\n  \"test\": \"test\"\n")
}

func TestMarshalReturnsError(t *testing.T) {
	_, err := json.Marshal(func() {})

	require.Error(t, err)
}

func TestUnmarshal(t *testing.T) {
	var msg map[string]string

	require.NoError(t, json.Unmarshal([]byte("{\"test\":\"test\"}"), &msg))
	require.Equal(t, map[string]string{"test": "test"}, msg)
}

func TestUnmarshalRejectsTrailingData(t *testing.T) {
	var msg map[string]string

	err := json.Unmarshal([]byte("{\"test\":\"test\"}{\"extra\":\"value\"}"), &msg)

	require.ErrorIs(t, err, errors.ErrTrailingData)
}

type message struct {
	Test string `json:"test"`
}
