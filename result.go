package ro2dtree

import (
	"sync/atomic"
)

//todo stats
type ResultPool struct {
	misses   int64
	capacity int
	list     chan *Result
}

func newResultPool(count, capacity int) *ResultPool {
	pool := &ResultPool{
		capacity: capacity,
		list:     make(chan *Result, count),
	}
	for i := 0; i < cap(pool.list); i++ {
		pool.list <- newResult(pool, capacity)
	}
	return pool
}

func (pool *ResultPool) Checkout() *Result {
	select {
	case Result := <-pool.list:
		return Result
	default:
		atomic.AddInt64(&pool.misses, 1)
		return newResult(nil, pool.capacity)
	}
}

type Result struct {
	target   Point
	position int
	polygons Polygons
	pool     *ResultPool
}

func newResult(pool *ResultPool, capacity int) *Result {
	return &Result{
		pool:     pool,
		polygons: make(Polygons, capacity),
	}
}

func (r *Result) Add(polygon Polygon) bool {
	r.polygons[r.position] = polygon
	r.position++
	return r.position != len(r.polygons)
}

func (r *Result) Polygons() Polygons {
	return r.polygons[:r.position]
}

func (r *Result) Close() {
	if r.pool != nil {
		r.position = 0
		r.pool.list <- r
	}
}

func (r *Result) Len() int {
	return r.position
}

func (r *Result) Less(i, j int) bool {
	return r.polygons[i].Centroid().DistanceTo(r.target) < r.polygons[j].Centroid().DistanceTo(r.target)
}

func (r *Result) Swap(i, j int) {
	r.polygons[i], r.polygons[j] = r.polygons[j], r.polygons[i]
}
