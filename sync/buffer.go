package sync

import "bytes"

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
func (p *BufferPool) Put(b *bytes.Buffer) {
	b.Reset()
	p.pool.Put(b)
}

// Copy the buffer to a []byte.
func (p *BufferPool) Copy(b *bytes.Buffer) []byte {
	buf := b.Bytes()
	newBuf := make([]byte, len(buf))

	copy(newBuf, buf)

	return newBuf
}
