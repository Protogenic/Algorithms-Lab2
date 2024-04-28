package solution

import (
	"Algorithms_Lab2/types"
	"fmt"
	"time"
)

func BruteForce(rectangles []types.Rectangle, points []types.Point) {
	startTime := time.Now()
	appearances := make([]int, len(points))
	for i := 0; i < len(points); i++ {
		for j := 0; j < len(rectangles); j++ {
			if rectangles[j].Belongs(points[i]) {
				appearances[i]++
			}
		}
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Println("Execution time: ", duration, "\n")
	//fmt.Println(appearances)
}
