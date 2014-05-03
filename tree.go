package ro2dtree

import (
	"math"
	"sort"
)

type Tree struct {
	min         int
	max         int
	Root        Polygon
	resultPool  *ResultPool
	scratchPool *ResultPool
}

func New(min, max int) *Tree {
	return &Tree{
		min:         min,
		max:         max,
		Root:        nil,
		resultPool:  newResultPool(32, 50),
		scratchPool: newResultPool(1024, max),
	}
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
	if t.Root.Contains(point) == false {
		return &Result{polygons: make(Polygons, 0)}
	}
	result := t.resultPool.Checkout()
	result.target = point
	t.find(t.Root, point, result)
	sort.Sort(result)
	return result
}

func (t *Tree) find(node Polygon, point Point, results *Result) bool {
	scratch := t.scratchPool.Checkout()
	defer scratch.Close()
	scratch.target = point

	children := node.Children()
	if children == nil {
		return results.Add(node)
	}
	for _, child := range children {
		if child.Contains(point) {
			scratch.Add(child)
		}
	}
	l := scratch.Len()
	if l == 0 {
		return true
	}
	if l != 1 {
		sort.Sort(scratch)
	}
	for _, child := range scratch.Polygons() {
		if t.find(child, point, results) == false {
			return false
		}
	}
	return true
}
