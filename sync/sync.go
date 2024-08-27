package sync

import (
	"bytes"
	"sync"
)

// BufferPool for sync.
type BufferPool struct {
	pool *sync.Pool
}

// NewBufferPool for sync.
func NewBufferPool() *BufferPool {
	pool := &sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}

	return &BufferPool{pool: pool}
}

// Get a new buffer.
func (p *BufferPool) Get() *bytes.Buffer {
	return p.pool.Get().(*bytes.Buffer)
}

// Put the buffer back.
func (p *BufferPool) Put(b *bytes.Buffer) {
	b.Reset()
	p.pool.Put(b)
}
