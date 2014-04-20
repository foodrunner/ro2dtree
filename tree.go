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
	t.Root = t.load(polygons)
}

func (t *Tree) load(polygons Polygons) Polygon {
	N := len(polygons)
	if N == 1 {
		return polygons[0]
	}
	sort.Sort(polygons)
	M := int(math.Ceil(float64(N / t.min)))
	parents := make(Polygons, 0, M)
	for i := 0; i < N; {
		node := &Node{children: make(Polygons, 0, t.max)}
		for j := 0; j < t.max; j++ {
			polygon := polygons[i]
			if j < t.min || node.IsClose(polygon) {
				node.Add(polygon)
			} else {
				break
				// this attempts to scan ahead for possible matches
				// but it turns loading into an O(N!) (I think) !!! and it actually
				// doesn't see to improve the overlap any..which I find impossible to believe

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
		node.seal()
		parents = append(parents, node)
	}
	return t.load(parents)
}

func (t *Tree) Find(point Point) *Result {
	//todo move this to a pool
	//todo allow offset / limit
	//enforce limit
	result := &Result{
		polygons: make(Polygons, 100000),
	}
	t.find(t.Root, point, result)
	return result
}

func (t *Tree) find(node Polygon, point Point, results *Result) {
	if node.Contains(point) {
		children := node.Children()
		if children == nil {
			results.Add(node)
		}
		for _, child := range children {
			t.find(child, point, results)
		}
	}
}
