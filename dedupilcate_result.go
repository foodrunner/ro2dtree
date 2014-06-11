package ro2dtree

type DeduplicateResult struct {
	position int
	items    Items
	pool     *ResultPool
	groupMap map[int]int
	target   Point
}

func DeduplicateResultFactory(pool *ResultPool, capacity int) Result {
	return &DeduplicateResult{
		position: 0,
		items:    make(Items, capacity),
		pool:     pool,
		groupMap: make(map[int]int),
	}
}

func (r *DeduplicateResult) Add(polygon Polygon) bool {
	rank := polygon.Centroid().DistanceTo(r.target)
	item := NewItem(polygon, rank)
	oldPosition, present := r.groupMap[item.Polygon().GroupId()]
	if present {
		oldItem := r.items[oldPosition]
		if oldItem.Rank() > item.Rank() {
			r.items[oldPosition] = item
		}
	} else {
		r.groupMap[item.Polygon().GroupId()] = r.position
		r.items[r.position] = item
		r.position++
	}
	return r.position != len(r.items)
}

func (r *DeduplicateResult) Items() Items {
	return r.items[:r.position]
}

func (r *DeduplicateResult) Close() {
	if r.pool != nil {
		r.position = 0
		r.groupMap = make(map[int]int)
		r.pool.list <- r
	}
}

func (r *DeduplicateResult) SetTarget(target Point) {
	r.target = target
}

func (r *DeduplicateResult) Len() int {
	return r.position
}

func (r *DeduplicateResult) Less(i, j int) bool {
	return r.items[i].Rank() < r.items[j].Rank()
}

func (r *DeduplicateResult) Swap(i, j int) {
	r.items[i], r.items[j] = r.items[j], r.items[i]
}
