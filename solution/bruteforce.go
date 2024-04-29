package solution

import (
	"Algorithms_Lab2/types"
	"fmt"
	"time"
)

func BruteForce(rectangles []types.Rectangle, points []types.Point) []int {
	startTime := time.Now()
	result := make([]int, len(points))
	for i := 0; i < len(points); i++ {
		for j := 0; j < len(rectangles); j++ {
			if rectangles[j].Belongs(points[i]) {
				result[i]++
			}
		}
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Println("Execution time: ", duration, "\n")

	return result
}
