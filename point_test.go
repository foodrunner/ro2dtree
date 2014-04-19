package ro2dtree

import (
	"sort"
	"testing"
)

func TestPointsGetSortedWithLocality(t *testing.T) {
	points := Points{
		NewPoint(100, 110),
		NewPoint(0, 2),
		NewPoint(1000, 2000),
		NewPoint(101, 120),
		NewPoint(10, 7),
		NewPoint(1500, 2200),
	}
	expected := Points{
		points[1],
		points[4],
		points[0],
		points[3],
		points[2],
		points[5],
	}
	sort.Sort(points)
	for index, point := range expected {
		if point != points[index] {
			t.Errorf("Expected point %d to be %v got %v", index, point, points[index])
		}
	}
}
