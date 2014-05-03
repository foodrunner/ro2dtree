package ro2dtree

import (
	"math"
	"sort"
)

type Tree struct {
	min        int
	max        int
	Root       Polygon
	stackPool  *StackPool
	resultPool *ResultPool
}

func New(minFill, maxFill, maxResults int) *Tree {
	maxDepth := 128 //TODO FIX
	return &Tree{
		min:        minFill,
		max:        maxFill,
		Root:       nil,
		stackPool:  newStackPool(32, maxDepth),
		resultPool: newResultPool(32, maxResults),
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
	if t.Root.Contains(point) == false {
		return EmptyResult
	}
	return t.find(t.Root, point)
}

func (t *Tree) find(node Polygon, point Point) *Result {
	stack := t.stackPool.Checkout()
	defer stack.Close()
	result := t.resultPool.Checkout()
	result.target = point

	for ; node != nil; node = stack.Pop() {
		if node.Contains(point) == false {
			continue
		}
		children := node.Children()
		if children == nil {
			if result.Add(node) == false {
				break
			}
			continue
		}
		for _, child := range children {
			stack.Push(child)
		}
	}
	return result
}
