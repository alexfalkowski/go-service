package test

import "github.com/alexfalkowski/go-service/v2/io"

// ErrReaderCloser is an io.ReadCloser test double whose methods fail with ErrFailed.
type ErrReaderCloser struct{}

// Read returns ErrFailed.
func (r *ErrReaderCloser) Read(_ []byte) (int, error) {
	return 0, ErrFailed
}

// Close returns ErrFailed.
func (r *ErrReaderCloser) Close() error {
	return ErrFailed
}

// StaticReader is an io.Reader test double that reads from Data.
type StaticReader struct {
	Data []byte
}

// Read copies Data into p and returns EOF if p was not filled.
func (r StaticReader) Read(p []byte) (int, error) {
	n := copy(p, r.Data)
	if n < len(p) {
		return n, io.EOF
	}

	return n, nil
}

// TrackingReader is an io.Reader test double that records read attempts.
type TrackingReader struct {
	Err   error
	Reads int
}

// Read increments Reads and returns Err.
func (r *TrackingReader) Read(_ []byte) (int, error) {
	r.Reads++
	return 0, r.Err
}
