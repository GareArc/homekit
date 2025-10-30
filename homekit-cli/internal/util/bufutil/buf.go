package bufutil

import (
	"bytes"
	"sync"
)

// Pool manages reusable buffers with size constraints
type Pool struct {
	pool        sync.Pool
	maxCapacity int // Maximum capacity to keep in pool (bytes)
}

const (
	KB = 1024
	MB = 1024 * KB
)

// NewPool creates a buffer pool
// initialSize: Starting buffer size (bytes)
// maxKeepSize: Max buffer size to keep in pool (larger buffers are discarded)
func NewPool(initialSize, maxKeepSize int) *Pool {
	return &Pool{
		maxCapacity: maxKeepSize,
		pool: sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 0, initialSize))
			},
		},
	}
}

// Get retrieves a buffer from the pool
func (p *Pool) Get() *bytes.Buffer {
	return p.pool.Get().(*bytes.Buffer)
}

// Put returns a buffer to the pool if it's within size limits
func (p *Pool) Put(buf *bytes.Buffer) {
	if buf.Cap() <= p.maxCapacity {
		buf.Reset()
		p.pool.Put(buf)
	}
	// Else: Buffer is too large, let it get GC'd
}

// WithBuffer executes fn with a buffer and automatically returns it
func (p *Pool) WithBuffer(fn func(*bytes.Buffer) error) error {
	buf := p.Get()
	defer p.Put(buf)
	return fn(buf)
}
