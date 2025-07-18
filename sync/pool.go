package sync

import "sync"

// NewPool of type T.
func NewPool[T any]() *Pool[T] {
	pool := &sync.Pool{
		New: func() any {
			return new(T)
		},
	}
	return &Pool[T]{pool: pool}
}

// Pool of type T.
type Pool[T any] struct {
	pool *sync.Pool
}

// Get an item of type T.
func (p *Pool[T]) Get() *T {
	return p.pool.Get().(*T)
}

// Put an item of type T back.
func (p *Pool[T]) Put(b *T) {
	p.pool.Put(b)
}
