package ro2dtree

//todo pool these bad boys

type Result struct {
	position int
	polygons Polygons
}

func (r *Result) Add(polygon Polygon) {
	r.polygons[r.position] = polygon
	r.position++
}

func (r *Result) Polygons() Polygons {
	return r.polygons[:r.position]
}
