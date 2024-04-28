package main

type Point struct {
	x int
	y int
}

func NewPoint(x int, y int) *Point {
	return &Point{x: x, y: y}
}

type Rectangle struct {
	lowerLeft  Point
	upperRight Point
}

func NewRectangle(x1, y1, x2, y2 int) *Rectangle {
	return &Rectangle{lowerLeft: *NewPoint(x1, y1), upperRight: *NewPoint(x2, y2)}
}

func (rectangle Rectangle) belongs(point Point) bool {
	return point.x >= rectangle.lowerLeft.x && point.x <= rectangle.upperRight.x &&
		point.y >= rectangle.lowerLeft.y && point.y <= rectangle.upperRight.y
}
