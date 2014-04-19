package ro2dtree

type Node struct {
	box      *Box
	Children Polygons
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
	n.Children = append(n.Children, polygon)
}

func (n *Node) Box() *Box {
	if n.box == nil {
		n.box = n.Children.Box()
	}
	return n.box
}
