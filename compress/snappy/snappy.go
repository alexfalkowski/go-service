package snappy

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	compress "github.com/alexfalkowski/go-service/v2/compress/errors"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/klauspost/compress/snappy"
)

// NewCompressor constructs a Snappy compressor implementation.
//
// The returned value implements [github.com/alexfalkowski/go-service/v2/compress.Compressor].
func NewCompressor() *Compressor {
	return &Compressor{}
}

// Compressor implements Snappy compression.
//
// It satisfies the [github.com/alexfalkowski/go-service/v2/compress.Compressor] interface.
type Compressor struct{}

// Compress returns the Snappy-compressed representation of data.
//
// An error is returned if data exceeds size.
func (c *Compressor) Compress(data []byte, size bytes.Size) ([]byte, error) {
	if int64(len(data)) > size.Bytes() {
		return nil, compress.ErrTooLarge
	}
	if snappy.MaxEncodedLen(len(data)) < 0 {
		return nil, compress.ErrTooLarge
	}

	return snappy.Encode(nil, data), nil
}

// Decompress returns the decompressed representation of data.
//
// An error is returned if data is not valid Snappy-encoded content or the decompressed data exceeds size.
func (c *Compressor) Decompress(data []byte, size bytes.Size) ([]byte, error) {
	decodedLen, err := snappy.DecodedLen(data)
	if err != nil {
		if errors.Is(err, snappy.ErrTooLarge) {
			return nil, compress.ErrTooLarge
		}

		return nil, err
	}
	if int64(decodedLen) > size.Bytes() {
		return nil, compress.ErrTooLarge
	}

	return snappy.Decode(nil, data)
}
