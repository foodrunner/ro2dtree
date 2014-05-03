package main

import (
	"fmt"
	"github.com/foodrunner/ro2dtree"
	"math/rand"
	"time"
)

const (
	GRID_MIN = 0
	GRID_MAX = 500
)

func main() {
	rand.Seed(20)
	all := createPolygons(20)
	tree := ro2dtree.New(8, 16)
	needle := ro2dtree.NewPoint(228, 30)

	tree.Load(all)
	s1 := time.Now()
	result := tree.Find(needle)
	fmt.Println(time.Now().Sub(s1), len(result.Polygons()))
	// fmt.Println("Tree Find:")
	// for _, polygon := range result.Polygons() {
	// 	fmt.Println(polygon)
	// }
	// fmt.Println("")

	// fmt.Println("Scan Find:")
	// s2 := time.Now()
	// found := find(all, needle)
	// e2 := time.Now()
	// for _, polygon := range found {
	// 	fmt.Println(polygon)
	// }
	// fmt.Println("")

	// fmt.Println(e1.Sub(s1))
	// fmt.Println(e2.Sub(s2))

	// 	fmt.Println(`
	// <style>
	// div{position:absolute;border:1px solid black;opacity:0.2}
	// </style>
	// `)
	// 	draw(tree.Root, 1, true)
	// 	fmt.Println(`
	// <script>
	// var divs = document.getElementsByTagName('div')
	// for (i = 0; i < divs.length; i++) {
	// 	divs[i].style.backgroundColor = '#' + ((i*10)+300).toString(16);
	// }
	// </script>
	// `)
}

func createPolygons(count int) ro2dtree.Polygons {
	polygons := make(ro2dtree.Polygons, count)
	for i := 0; i < count; i++ {
		polygons[i] = createPolygon()
	}
	return polygons
}

func createPolygon() ro2dtree.Polygon {
	lengthA := float64(rand.Int31n(5) + 50)
	lengthB := float64(rand.Int31n(5) + 50)
	x := float64(rand.Int31n(GRID_MAX - int32(lengthA)))
	y := float64(rand.Int31n(GRID_MAX - int32(lengthB)))

	return ro2dtree.NewPolygon(ro2dtree.Points{
		ro2dtree.NewPoint(x, y),
		ro2dtree.NewPoint(x+lengthA, y),
		ro2dtree.NewPoint(x+lengthA, y+lengthB),
		ro2dtree.NewPoint(x, y+lengthB),
		ro2dtree.NewPoint(x, y),
	})
}

func draw(polygon ro2dtree.Polygon, level int, recurse bool) {
	var box *ro2dtree.Box

	if node, ok := polygon.(*ro2dtree.Node); ok {
		box = node.Box()
	} else {
		box = polygon.(*ro2dtree.SimplePolygon).Box()
		recurse = false
	}
	topLeft, bottomRight := box.TopLeft, box.BottomRight
	fmt.Println(fmt.Sprintf(`<div style="left:%dpx;top:%dpx;width:%dpx;height:%dpx"></div>`, int(topLeft.X), int(topLeft.Y), int(bottomRight.X-topLeft.X), int(bottomRight.Y-topLeft.Y)))
	if recurse {
		for _, child := range polygon.Children() {
			draw(child, level+1, false)
		}
	}
}

func find(polygons ro2dtree.Polygons, needle ro2dtree.Point) ro2dtree.Polygons {
	results := make(ro2dtree.Polygons, 0, 1000)
	for _, polygon := range polygons {
		if polygon.Contains(needle) {
			results = append(results, polygon)
		}
	}
	return results
}
