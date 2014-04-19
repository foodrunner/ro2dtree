package ro2dtree

type Point struct {
	X float64
	Y float64
}

func (p Point) Id() uint64 {
	//adjust for negatives
	x := uint32((p.X + 180) * 1000000)
	y := uint32((p.Y + 90) * 1000000)
	var result uint64
	for i := uint16(31); i > 0; i-- {
		result |= uint64((x >> i) & 1)
		result = result << 1
		result |= uint64((y >> i) & 1)
		result = result << 1
	}
	result |= uint64((x >> 0) & 1)
	result |= uint64((y >> 0) & 1)
	return result
}

func NewPoint(x, y float64) Point {
	return Point{x, y}
}

type Points []Point

func (p Points) Len() int {
	return len(p)
}

func (p Points) Less(i, j int) bool {
	return p[i].Id() < p[j].Id()
}

func (p Points) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
