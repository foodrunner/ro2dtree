package ro2dtree

var EmptyResult = NormalResultFactory(nil, 0)

type NormalResult struct {
	position int
	items    Items
	pool     *ResultPool
}

func NormalResultFactory(pool *ResultPool, capacity int) Result {
	return &NormalResult{
		position: 0,
		items:    make(Items, capacity),
		pool:     pool,
	}
}

func (r *NormalResult) Add(item *Item) bool {
	r.items[r.position] = item
	r.position++
	return r.position != len(r.items)
}

func (r *NormalResult) Items() Items {
	return r.items[:r.position]
}

func (r *NormalResult) Close() {
	if r.pool != nil {
		r.position = 0
		r.pool.list <- r
	}
}

func (r *NormalResult) Len() int {
	return r.position
}

func (r *NormalResult) Less(i, j int) bool {
	return r.items[i].Rank() < r.items[j].Rank()
}

func (r *NormalResult) Swap(i, j int) {
	r.items[i], r.items[j] = r.items[j], r.items[i]
}
