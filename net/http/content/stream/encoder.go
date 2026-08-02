package stream

import "github.com/alexfalkowski/go-service/v2/encoding/stream"

// requestEncoder keeps bidirectional stream finalization on the response stream while closing both codecs.
type requestEncoder struct {
	stream.Encoder
	decoder stream.Decoder
}

func (e *requestEncoder) Close() error {
	if err := e.Encoder.Close(); err != nil {
		_ = e.decoder.Close()

		return err
	}

	return e.decoder.Close()
}
