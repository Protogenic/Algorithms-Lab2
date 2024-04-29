package main

import (
	"Algorithms_Lab2/solution"
	"Algorithms_Lab2/types"
	"fmt"
	"math/rand"
)

func main() {
	//rectangles, points := readData()
	rectangles, points := randomData()
	fmt.Println("Brute force")
	solution.BruteForce(rectangles, points)
	//fmt.Println("Compressed map")
	//solution.CompressedMap(rectangles, points)
	fmt.Println("Segment tree")
	solution.SegmentTree(rectangles, points)
}

func randomData() ([]types.Rectangle, []types.Point) {
	var numberOfRectangles, numberOfPoints int
	fmt.Scanf("%d %d", &numberOfRectangles, &numberOfPoints)
	max := int(1000000000)

	rectangles := make([]types.Rectangle, numberOfRectangles)
	for i := 0; i < numberOfRectangles; i++ {
		var x1, y1, x2, y2 int
		x1 = rand.Intn(max)
		x2 = rand.Intn(max-x1) + x1 + 1
		y1 = rand.Intn(max)
		y2 = rand.Intn(max-y1) + y1 + 1
		rectangles[i] = *types.NewRectangle(x1, y1, x2, y2)
	}

	points := make([]types.Point, numberOfPoints)
	for i := 0; i < numberOfPoints; i++ {
		var x, y int
		x = rand.Intn(max)
		y = rand.Intn(max)
		points[i] = *types.NewPoint(x, y)
	}

	return rectangles, points
}

func readData() ([]types.Rectangle, []types.Point) {
	var numberOfRectangles int
	fmt.Scanf("%d\n", &numberOfRectangles)
	rectangles := make([]types.Rectangle, numberOfRectangles)
	for i := 0; i < numberOfRectangles; i++ {
		var x1, y1, x2, y2 int
		fmt.Scanf("%d %d %d %d\n", &x1, &y1, &x2, &y2)
		rectangles[i] = *types.NewRectangle(x1, y1, x2, y2)
	}

	var numberOfPoints int
	fmt.Scanf("%d\n", &numberOfPoints)
	points := make([]types.Point, numberOfPoints)
	for i := 0; i < numberOfPoints; i++ {
		var x, y int
		fmt.Scanf("%d %d\n", &x, &y)
		points[i] = *types.NewPoint(x, y)
	}

	return rectangles, points
}
