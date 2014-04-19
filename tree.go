package ro2dtree

import (
	"math"
	"sort"
)

type Tree struct {
	min  int
	max  int
	Root Polygon
}

func New(min, max int) *Tree {
	return &Tree{min, max, nil}
}

func (t *Tree) Load(polygons Polygons) {
	t.Root = t.build(polygons)
}

func (t *Tree) build(polygons Polygons) Polygon {
	N := len(polygons)
	if N == 1 {
		return polygons[0]
	}
	sort.Sort(polygons)
	M := int(math.Ceil(float64(N / t.min)))
	parents := make(Polygons, 0, M)
	for i := 0; i < N; {
		node := &Node{Children: make(Polygons, 0, t.max)}
		for j := 0; j < t.max; j++ {
			polygon := polygons[i]
			if j < t.min || node.IsClose(polygon) {
				node.Add(polygon)
			} else {
				break
				// next := true
				// for k := i+1; k < N; k++ {
				// 	polygon = polygons[k]
				// 	if node.IsClose(polygon) {
				// 		polygons[k-1], polygons[k] = polygons[k], polygons[k - 1]
				// 		next = false
				// 	}
				// }
				// if next {
				// 	break
				// } else {
				// 	j--
				// 	i--
				// }
			}
			i++
			if i == N {
				break
			}
		}
		parents = append(parents, trim(node))
	}
	return t.build(parents)
}

func trim(node *Node) *Node {
	if l := len(node.Children); l < cap(node.Children) {
		trimmed := make(Polygons, l)
		copy(trimmed, node.Children)
		node.Children = trimmed
	}
	return node
}
