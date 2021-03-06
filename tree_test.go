package ro2dtree

import (
	"math/rand"
	"testing"
)

func TestTreeFindMatchingValues(t *testing.T) {
	for i := 0; i < 100; i++ {
		tree := New(2, 4, 1000, NormalResultFactory)
		polygons := createPolygons(500)
		tree.Load(polygons)
		needle := NewPoint(float64(rand.Int31n(20)), float64(rand.Int31n(20)))
		expectSameNodes(t, tree.Find(needle).Items(), scan(polygons, needle))
	}
}

func TestTreeFindUniqueByGroupMatchingValues(t *testing.T) {
	polygons := make(Polygons, 4)
	points0 := Points{
		NewPoint(0, 0),
		NewPoint(0, 4),
		NewPoint(4, 4),
		NewPoint(4, 0),
		NewPoint(0, 0),
	}
	polygons[0] = NewPolygon(0, 1, points0)

	points1 := Points{
		NewPoint(0, 0),
		NewPoint(0, 2),
		NewPoint(2, 2),
		NewPoint(2, 0),
		NewPoint(0, 0),
	}
	group1Polygon := NewPolygon(1, 1, points1)
	polygons[1] = group1Polygon

	points2 := Points{
		NewPoint(0, 0),
		NewPoint(0, 2),
		NewPoint(4, 2),
		NewPoint(4, 0),
		NewPoint(0, 0),
	}
	polygons[2] = NewPolygon(2, 1, points2)

	points3 := Points{
		NewPoint(0, 0),
		NewPoint(0, 4),
		NewPoint(2, 4),
		NewPoint(2, 0),
		NewPoint(0, 0),
	}
	group2Polygon := NewPolygon(3, 2, points3)
	polygons[3] = group2Polygon

	tree := New(2, 4, 100, DeduplicateResultFactory)
	tree.Load(polygons)
	needle := NewPoint(1, 1)
	itemsFound := tree.Find(needle).Items()
	// Load expected polygons
	expected := make(map[Polygon]struct{})
	expected[group1Polygon] = struct{}{}
	expected[group2Polygon] = struct{}{}
	expectSameNodes(t, itemsFound, expected)
}

func TestTreeGetPolygonById(t *testing.T) {
	tree := New(8, 16, 1000, NormalResultFactory)
	polygons := createPolygons(500)
	tree.Load(polygons)
	idToFind := 10
	p := tree.Get(idToFind)
	if p.Id() != idToFind {
		t.Errorf("Get polygon with id %v but found %v", idToFind, p.Id())
	}
	idOutRange := 10000000
	p = tree.Get(idOutRange)
	if p != nil {
		t.Errorf("Tree should return nil when get id %v", idOutRange)
	}
}

func TestHitTest(t *testing.T) {
	tree := New(2, 4, 100, NormalResultFactory)
	polygons := make(Polygons, 4)

	points0 := Points{
		NewPoint(0, 0),
		NewPoint(0, 2),
		NewPoint(2, 2),
		NewPoint(2, 0),
		NewPoint(0, 0),
	}
	polygons[0] = NewPolygon(0, 1, points0)

	points1 := Points{
		NewPoint(0, 0),
		NewPoint(0, 4),
		NewPoint(4, 4),
		NewPoint(4, 0),
		NewPoint(0, 0),
	}
	polygons[1] = NewPolygon(1, 1, points1)

	points2 := Points{
		NewPoint(2, 2),
		NewPoint(2, 4),
		NewPoint(4, 4),
		NewPoint(4, 2),
		NewPoint(2, 2),
	}
	polygons[2] = NewPolygon(2, 1, points2)

	points3 := Points{
		NewPoint(10, 2),
		NewPoint(10, 6),
		NewPoint(14, 6),
		NewPoint(14, 2),
		NewPoint(10, 2),
	}
	polygons[3] = NewPolygon(3, 1, points3)

	tree.Load(polygons)

	point := NewPoint(3, 3)

	ids := []int{0, 1, 2, 3}
	id := tree.HitTest(ids, point)
	if id != 2 {
		t.Errorf("Expecting hit test %v, get %v", 2, id)
	}

	ids = []int{0, 3}
	id = tree.HitTest(ids, point)
	if id != -1 {
		t.Errorf("Expecting hit test fails, get %v", id)
	}
}

func BenchmarkTreeFindLowFill(b *testing.B) {
	tree := New(2, 4, 1000, NormalResultFactory)
	polygons := createPolygons(50000)
	tree.Load(polygons)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		needle := NewPoint(float64(rand.Int31n(20)), float64(rand.Int31n(20)))
		tree.Find(needle)
	}
}

func BenchmarkTreeFindHighFill(b *testing.B) {
	tree := New(8, 16, 1000, NormalResultFactory)
	polygons := createPolygons(50000)
	tree.Load(polygons)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		needle := NewPoint(float64(rand.Int31n(20)), float64(rand.Int31n(20)))
		tree.Find(needle)
	}
}

func expectSameNodes(t *testing.T, actual Items, expected map[Polygon]struct{}) {
	if len(actual) != len(expected) {
		t.Errorf("Expecting %d results got %d", len(expected), len(actual))
	}
	for index, item := range actual {
		polygon := item.Polygon()
		if _, exists := expected[polygon]; !exists {
			t.Errorf("Polygon %v at index %d should not have been found", polygon, index)
		}
	}
}

func createPolygons(count int) Polygons {
	polygons := make(Polygons, count)
	for i := 0; i < count; i++ {
		polygons[i] = createPolygon(i)
	}
	return polygons
}

func createPolygon(id int) Polygon {
	lengthA := float64(rand.Int31n(50) + 10)
	lengthB := float64(rand.Int31n(50) + 10)
	x := float64(rand.Int31n(120 - int32(lengthA)))
	y := float64(rand.Int31n(120 - int32(lengthB)))

	return NewPolygon(id, 1, Points{
		NewPoint(x, y),
		NewPoint(x+lengthA, y),
		NewPoint(x+lengthA, y+lengthB),
		NewPoint(x, y+lengthB),
		NewPoint(x, y),
	})
}

func scan(haystack Polygons, needle Point) map[Polygon]struct{} {
	results := make(map[Polygon]struct{})
	for _, polygon := range haystack {
		if polygon.Contains(needle) {
			results[polygon] = struct{}{}
		}
	}
	return results
}
