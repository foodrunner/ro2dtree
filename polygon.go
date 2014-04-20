package ro2dtree

type Box struct {
	TopLeft     Point
	BottomRight Point
}

type Polygon interface {
	Box() *Box
	Centroid() Point
	Contains(p Point) bool
	Children() Polygons
}

type Polygons []Polygon

func (p Polygons) Len() int {
	return len(p)
}

func (p Polygons) Less(i, j int) bool {
	return p[i].Centroid().Id() < p[j].Centroid().Id()
}

func (p Polygons) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p Polygons) Box() *Box {
	box := p[0].Box()
	topLeft, bottomRight := box.TopLeft, box.BottomRight
	for i, l := 1, len(p); i < l; i++ {
		box := p[i].Box()
		if box.TopLeft.X < topLeft.X {
			topLeft.X = box.TopLeft.X
		}
		if box.BottomRight.X > bottomRight.X {
			bottomRight.X = box.BottomRight.X
		}
		if box.TopLeft.Y < topLeft.Y {
			topLeft.Y = box.TopLeft.Y
		}
		if box.BottomRight.Y > bottomRight.Y {
			bottomRight.Y = box.BottomRight.Y
		}
	}
	return &Box{topLeft, bottomRight}
}

// func (p Polygons) IsClose(polygon Polygon) bool {
// 	topLeft, bottomRight := p.BoundingBox()
// 	//overlap
// 	if topLeft.X < polygon.BottomRight.X && bottomRight.X > polygon.TopLeft.X && topLeft.Y < polygon.BottomRight.Y && bottomRight.Y > polygon.TopLeft.Y {
// 		return true
// 	}
// 	return false
// }

type SimplePolygon struct {
	points   Points
	box      *Box
	centroid Point
}

func NewPolygon(points Points) *SimplePolygon {
	return &SimplePolygon{
		points:   points,
		box:      calculateBox(points),
		centroid: calculateCentroid(points),
	}
}

func (p *SimplePolygon) Box() *Box {
	return p.box
}

func (p *SimplePolygon) Centroid() Point {
	return p.centroid
}

func (p *SimplePolygon) Children() Polygons {
	return nil
}

func (p *SimplePolygon) Contains(point Point) bool {
	topLeft, bottomRight := p.box.TopLeft, p.box.BottomRight
	x, y := point.X, point.Y
	if x <= topLeft.X || y <= topLeft.Y || x >= bottomRight.X || y >= bottomRight.Y {
		return false
	}
	var hit bool
	points := p.points
	l := len(points)
	for i, j := 0, l-1; i < l; {
		if ((points[i].X > x) != (points[j].X > y)) && (x < (points[j].Y-points[i].Y)*(y-points[i].X)/(points[j].X-points[i].X)+points[i].Y) {
			hit = !hit
		}
		j = i
		i++
	}
	return hit
}

func calculateBox(points Points) *Box {
	first := points[0]
	x1, y1, x2, y2 := first.X, first.Y, first.X, first.Y
	// last point is the same as the first
	for i, l := 1, len(points); i < l-1; i++ {
		point := points[i]
		if point.X < x1 {
			x1 = point.X
		} else if point.X > x2 {
			x2 = point.X
		}
		if point.Y < y1 {
			y1 = point.Y
		} else if point.Y > y2 {
			y2 = point.Y
		}
	}
	return &Box{NewPoint(x1, y1), NewPoint(x2, y2)}
}

func calculateCentroid(points Points) Point {
	l := len(points)
	x, y, area, f := 0.0, 0.0, 0.0, 0.0
	for i, j := 0, l-1; i < l; {
		p1, p2 := points[i], points[j]
		f = p1.X*p2.Y - p2.X*p1.Y
		x += (p1.X + p2.X) * f
		y += (p1.Y + p2.Y) * f
		area += f
		j = i
		i++
	}
	f = area * 3
	return NewPoint(x/f, y/f)
}
