package solution

import (
	"fmt"
	"sort"
	"time"
)

import (
	"Algorithms_Lab2/types"
)

func CompressedMap(rectangles []types.Rectangle, points []types.Point) {
	compressedX, compressedY, compressedMap := CompressedMapUtil(rectangles)

	startTime := time.Now()

	for i := range points {
		X := binarySearch(compressedX, points[i].X)
		Y := binarySearch(compressedY, points[i].Y)

		if X == -1 || Y == -1 {
			fmt.Println(points[i], 0)
		} else {
			fmt.Println(points[i], compressedMap[Y][X])
		}
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Println("Execution time: ", duration, "\n")
}

func CompressedMapUtil(rectangles []types.Rectangle) ([]int, []int, [][]int) {
	fmt.Print("Preparation time: ")
	startTime := time.Now()
	uniqueX := make(map[int]struct{})
	uniqueY := make(map[int]struct{})

	for _, rectangle := range rectangles {
		uniqueX[rectangle.LowerLeft.X] = struct{}{}
		uniqueY[rectangle.LowerLeft.Y] = struct{}{}

		uniqueX[rectangle.UpperRight.X+1] = struct{}{}
		uniqueY[rectangle.UpperRight.Y+1] = struct{}{}
	}

	var compressedX = make([]int, len(uniqueX))
	for i := range uniqueX {
		compressedX = append(compressedX, i)
	}
	sort.Slice(compressedX, func(i, j int) bool {
		return compressedX[i] < compressedX[j]
	})

	var compressedY = make([]int, len(uniqueY))
	for i := range uniqueY {
		compressedY = append(compressedY, i)
	}
	sort.Slice(compressedY, func(i, j int) bool {
		return compressedY[i] < compressedY[j]
	})

	compressedMap := make([][]int, len(compressedY))
	for i := range compressedMap {
		compressedMap[i] = make([]int, len(compressedX))
	}

	for _, rectangle := range rectangles {
		leftX := binarySearch(compressedX, rectangle.LowerLeft.X)
		leftY := binarySearch(compressedY, rectangle.LowerLeft.Y)

		rightX := binarySearch(compressedX, rectangle.UpperRight.X+1)
		rightY := binarySearch(compressedY, rectangle.UpperRight.Y+1)

		for i := leftY; i < rightY; i++ {
			for j := leftX; j < rightX; j++ {
				compressedMap[i][j]++
			}
		}
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Println(duration)

	return compressedX, compressedY, compressedMap
}

func binarySearch(array []int, target int) int {
	left := 0
	right := len(array)

	for left < right {
		mid := (right + left) / 2
		if target >= array[mid] {
			left = mid + 1
		} else {
			right = mid
		}
	}

	return left - 1
}
