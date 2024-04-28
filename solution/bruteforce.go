package main

func BruteForce(rectangles []Rectangle, points []Point) []int {
	appearances := make([]int, len(points))
	for i := 0; i < len(points); i++ {
		for j := 0; j < len(rectangles); j++ {
			if rectangles[j].belongs(points[i]) {
				appearances[i]++
			}
		}
	}
	return appearances
}
