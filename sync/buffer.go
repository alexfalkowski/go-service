package sync

import (
	"sync"

	"github.com/alexfalkowski/go-service/v2/bytes"
)

// NewBufferPool for sync.
func NewBufferPool() *BufferPool {
	return &BufferPool{pool: &sync.Pool{
		New: func() any {
			return &bytes.Buffer{}
		},
	}}
}

// BufferPool for sync.
type BufferPool struct {
	pool *sync.Pool
}

// Get a new buffer.
func (p *BufferPool) Get() *bytes.Buffer {
	return p.pool.Get().(*bytes.Buffer)
}

// Put the buffer back.
func (p *BufferPool) Put(buffer *bytes.Buffer) {
	buffer.Reset()
	p.pool.Put(buffer)
}

// Copy the buffer to a []byte.
func (p *BufferPool) Copy(buffer *bytes.Buffer) []byte {
	return bytes.Clone(buffer.Bytes())
}
