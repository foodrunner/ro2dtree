package ro2dtree

type Node struct {
	box      *Box
	children Polygons
}

//overlap
func (n *Node) IsClose(polygon Polygon) bool {
	box := n.Box()
	target := polygon.Box()
	return box.TopLeft.X < target.BottomRight.X && box.BottomRight.X > target.TopLeft.X && box.TopLeft.Y < target.BottomRight.Y && box.BottomRight.Y > target.TopLeft.Y
}

func (n *Node) Centroid() Point {
	box := n.Box()
	return NewPoint((box.BottomRight.X-box.TopLeft.X)/2, (box.BottomRight.Y-box.TopLeft.Y)/2)
}

func (n *Node) Add(polygon Polygon) {
	n.box = nil
	n.children = append(n.children, polygon)
}

func (n *Node) Box() *Box {
	if n.box == nil {
		n.box = n.children.Box()
	}
	return n.box
}

func (n *Node) Children() Polygons {
	return n.children
}

func (n *Node) Contains(point Point) bool {
	topLeft, bottomRight := n.box.TopLeft, n.box.BottomRight
	x, y := point.X, point.Y
	return x >= topLeft.X && y >= topLeft.Y && x <= bottomRight.X && y <= bottomRight.Y
}

func (n *Node) String() string {
	return n.Centroid().String()
}

func (n *Node) Id() int {
	return -1
}

func (n *Node) GroupId() int {
	return -1
}

func (n *Node) seal() {
	if l := len(n.children); l < cap(n.children) {
		trimmed := make(Polygons, l)
		copy(trimmed, n.children)
		n.children = trimmed
	}
	//make sure n.box is valid
	n.Box()
}
