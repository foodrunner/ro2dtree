package ro2dtree

import (
	"math/rand"
	"testing"
)

func TestTreeFindMatchingValues(t *testing.T) {
	for i := 0; i < 100; i++ {
		tree := New(2, 4, 1000)
		polygons := createPolygons(500)
		tree.Load(polygons)
		needle := NewPoint(float64(rand.Int31n(20)), float64(rand.Int31n(20)))
		expectSameNodes(t, tree.Find(needle).Polygons(), scan(polygons, needle))
	}
}

func BenchmarkTreeFindLowFill(b *testing.B) {
	tree := New(2, 4, 1000)
	polygons := createPolygons(50000)
	tree.Load(polygons)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		needle := NewPoint(float64(rand.Int31n(20)), float64(rand.Int31n(20)))
		tree.Find(needle)
	}
}

func BenchmarkTreeFindHighFill(b *testing.B) {
	tree := New(8, 16, 1000)
	polygons := createPolygons(50000)
	tree.Load(polygons)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		needle := NewPoint(float64(rand.Int31n(20)), float64(rand.Int31n(20)))
	}
}

func expectSameNodes(t *testing.T, actual Polygons, expected map[Polygon]struct{}) {
	if len(actual) != len(expected) {
		t.Errorf("Expecting %d results got %d", len(expected), len(actual))
	}
	for index, polygon := range actual {
		if _, exists := expected[polygon]; !exists {
			t.Errorf("Polygon %v at index %d should not have been found", polygon, index)
		}
	}
}

func createPolygons(count int) Polygons {
	polygons := make(Polygons, count)
	for i := 0; i < count; i++ {
		polygons[i] = createPolygon()
	}
	return polygons
}

func createPolygon() Polygon {
	lengthA := float64(rand.Int31n(50) + 10)
	lengthB := float64(rand.Int31n(50) + 10)
	x := float64(rand.Int31n(120 - int32(lengthA)))
	y := float64(rand.Int31n(120 - int32(lengthB)))

	return NewPolygon(Points{
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
