package base64_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding/base64"
	"github.com/stretchr/testify/require"
)

func TestEncodeUsesStandardPaddedBase64(t *testing.T) {
	tests := []struct {
		expected string
		name     string
		src      []byte
	}{
		{name: "empty", src: nil, expected: ""},
		{name: "one byte", src: []byte("f"), expected: "Zg=="},
		{name: "two bytes", src: []byte("fo"), expected: "Zm8="},
		{name: "three bytes", src: []byte("foo"), expected: "Zm9v"},
		{name: "standard alphabet", src: []byte{0xfb, 0xff}, expected: "+/8="},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, base64.Encode(tt.src))
		})
	}
}

func TestDecodeUsesStandardPaddedBase64(t *testing.T) {
	tests := []struct {
		name     string
		encoded  string
		expected []byte
	}{
		{name: "empty", encoded: "", expected: []byte{}},
		{name: "one byte", encoded: "Zg==", expected: []byte("f")},
		{name: "two bytes", encoded: "Zm8=", expected: []byte("fo")},
		{name: "three bytes", encoded: "Zm9v", expected: []byte("foo")},
		{name: "standard alphabet", encoded: "+/8=", expected: []byte{0xfb, 0xff}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := base64.Decode(tt.encoded)
			require.NoError(t, err)
			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestEncodedLenUsesStandardPadding(t *testing.T) {
	tests := []struct {
		name     string
		size     bytes.Size
		expected int64
	}{
		{name: "empty", size: 0, expected: 0},
		{name: "one byte", size: 1, expected: 4},
		{name: "two bytes", size: 2, expected: 4},
		{name: "three bytes", size: 3, expected: 4},
		{name: "four bytes", size: 4, expected: 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, base64.EncodedLen(tt.size))
		})
	}
}
