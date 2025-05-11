package sync

import "github.com/alexfalkowski/go-service/bytes"

// NewBufferPool for sync.
func NewBufferPool() *BufferPool {
	pool := NewPool[bytes.Buffer]()

	return &BufferPool{pool: pool}
}

// BufferPool for sync.
type BufferPool struct {
	pool *Pool[bytes.Buffer]
}

// Get a new buffer.
func (p *BufferPool) Get() *bytes.Buffer {
	return p.pool.Get()
}

// Put the buffer back.
func (p *BufferPool) Put(buffer *bytes.Buffer) {
	buffer.Reset()
	p.pool.Put(buffer)
}

// Copy the buffer to a []byte.
func (p *BufferPool) Copy(buffer *bytes.Buffer) []byte {
	return bytes.Copy(buffer.Bytes())
}
