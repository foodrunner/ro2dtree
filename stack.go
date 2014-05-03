package ro2dtree

import (
	"sync/atomic"
)

type StackPool struct {
	misses   int64
	capacity int
	list     chan *Stack
}

func newStackPool(count, capacity int) *StackPool {
	pool := &StackPool{
		capacity: capacity,
		list:     make(chan *Stack, count),
	}
	for i := 0; i < cap(pool.list); i++ {
		pool.list <- newStack(pool, capacity)
	}
	return pool
}

func (pool *StackPool) Checkout() *Stack {
	select {
	case stack := <-pool.list:
		return stack
	default:
		atomic.AddInt64(&pool.misses, 1)
		return newStack(nil, pool.capacity)
	}
}

type Stack struct {
	pool     *StackPool
	position int
	polygons Polygons
}

func newStack(pool *StackPool, size int) *Stack {
	return &Stack{
		pool:     pool,
		polygons: make(Polygons, size),
	}
}

func (s *Stack) Push(polygon Polygon) {
	s.polygons[s.position] = polygon
	s.position++
}

func (s *Stack) Pop() Polygon {
	if s.position == 0 {
		return nil
	}
	s.position--
	return s.polygons[s.position]
}

func (s *Stack) Close() {
	if s.pool != nil {
		s.position = 0
		s.pool.list <- s
	}
}
