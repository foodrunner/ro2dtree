package ro2dtree

type Item struct {
	polygon Polygon
	rank    float64
}

func NewItem(polygon Polygon, rank float64) *Item {
	return &Item{
		polygon: polygon,
		rank:    rank,
	}
}

func (item *Item) Rank() float64 {
	return item.rank
}

func (item *Item) Polygon() Polygon {
	return item.polygon
}

type Items []*Item
