package syncx

import (
	"sync"
	"sync/atomic"
)

type PoolWrapper[T any] struct {
	unreclaimed int64
	p           sync.Pool
}

func NewPoolWrapper[T any](builder func() T) PoolWrapper[T] {
	return PoolWrapper[T]{
		p: sync.Pool{
			New: func() any { return builder() },
		},
	}
}

func (p *PoolWrapper[T]) Get() T {
	atomic.AddInt64(&p.unreclaimed, 1)
	return p.p.Get().(T)
}

func (p *PoolWrapper[T]) Put(o T) {
	atomic.AddInt64(&p.unreclaimed, -1)
	p.p.Put(o)
}

func (p *PoolWrapper[T]) UnreclaimedCount() int64 {
	return atomic.LoadInt64(&p.unreclaimed)
}
