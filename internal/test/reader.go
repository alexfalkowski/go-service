package test

// ErrReaderCloser for test.
type ErrReaderCloser struct{}

// Read returns ErrFailed.
func (r *ErrReaderCloser) Read(_ []byte) (int, error) {
	return 0, ErrFailed
}

// Close returns ErrFailed.
func (r *ErrReaderCloser) Close() error {
	return ErrFailed
}
