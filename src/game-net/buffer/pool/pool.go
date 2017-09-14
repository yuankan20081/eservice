package pool

import (
	"game-net/buffer"
	"sync"
)

type Pool struct {
	p *sync.Pool
}

func newBuffer() interface{} {
	return buffer.New()
}

func New() *Pool {
	return &Pool{
		p: &sync.Pool{
			New: newBuffer,
		},
	}
}

func (p *Pool) Get() *buffer.Buffer {
	return p.p.Get().(*buffer.Buffer)
}

func (p *Pool) Put(b *buffer.Buffer) {
	p.p.Put(b)
}
