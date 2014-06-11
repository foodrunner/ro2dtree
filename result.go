package ro2dtree

import (
	"sync/atomic"
)

type Result interface {
	Items() Items
	Add(item *Item) bool
	Close()
	Len() int
}

type ResultFactoryFunc func(*ResultPool, int) Result

//todo stats
type ResultPool struct {
	misses   int64
	capacity int
	list     chan Result
	factory  ResultFactoryFunc
}

func newResultPool(count, capacity int, fn ResultFactoryFunc) *ResultPool {
	pool := &ResultPool{
		capacity: capacity,
		list:     make(chan Result, count),
		factory:  fn,
	}
	for i := 0; i < cap(pool.list); i++ {
		pool.list <- fn(pool, capacity)
	}
	return pool
}

func (pool *ResultPool) Checkout() Result {
	select {
	case result := <-pool.list:
		return result
	default:
		atomic.AddInt64(&pool.misses, 1)
		return pool.factory(nil, pool.capacity)
	}
}
