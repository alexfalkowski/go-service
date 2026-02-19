package compress_test

import (
	"errors"
	"testing"

	"github.com/alexfalkowski/go-service/v2/compress"
	"github.com/alexfalkowski/go-service/v2/compress/none"
	"github.com/alexfalkowski/go-service/v2/compress/s2"
	"github.com/alexfalkowski/go-service/v2/compress/snappy"
	"github.com/alexfalkowski/go-service/v2/compress/zstd"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	kinds := []string{"zstd", "s2", "snappy", "none"}
	payloads := [][]byte{
		{},
		strings.Bytes("hello"),
		{0x00, 0x01, 0x02, 0xff, 0x00},
		makeLarge(1 << 16),
	}

	for _, kind := range kinds {
		t.Run(kind, func(t *testing.T) {
			cmp := test.Compressor.Get(kind)
			require.NotNil(t, cmp)

			for _, data := range payloads {
				d := cmp.Compress(data)

				ns, err := cmp.Decompress(d)
				require.NoError(t, err)
				requireRoundTrip(t, data, ns)
			}
		})
	}

	for _, key := range []string{"test", "bob"} {
		cmp := test.Compressor.Get(key)
		require.Nil(t, cmp)
	}
}

func TestMapRegister(t *testing.T) {
	m := compress.NewMap(compress.MapParams{
		Zstd:   zstd.NewCompressor(),
		S2:     s2.NewCompressor(),
		Snappy: snappy.NewCompressor(),
		None:   none.NewCompressor(),
	})

	original := m.Get("zstd")
	require.NotNil(t, original)

	custom := test.NewCompressor(errors.New("boom"))
	m.Register("zstd", custom)

	require.Same(t, custom, m.Get("zstd"))
	require.Nil(t, m.Get("missing"))
}

func TestDecompressCorruptData(t *testing.T) {
	cases := []struct {
		kind string
		data []byte
	}{
		{kind: "zstd", data: []byte{0xff, 0xff, 0xff}},
		{kind: "s2", data: []byte{0xff}},
		{kind: "snappy", data: []byte{0xff}},
	}

	for _, tc := range cases {
		t.Run(tc.kind, func(t *testing.T) {
			cmp := test.Compressor.Get(tc.kind)
			require.NotNil(t, cmp)

			_, err := cmp.Decompress(tc.data)
			require.Error(t, err)
		})
	}
}

func makeLarge(size int) []byte {
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(i)
	}
	return data
}

func requireRoundTrip(t *testing.T, original []byte, decoded []byte) {
	t.Helper()

	if len(original) == 0 {
		require.Empty(t, decoded)
		return
	}

	require.Equal(t, original, decoded)
}
