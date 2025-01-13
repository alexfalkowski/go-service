package test

// BadReader for test.
type BadReader struct{}

// Read returns ErrFailed.
func (r *BadReader) Read(_ []byte) (int, error) {
	return 0, ErrFailed
}
