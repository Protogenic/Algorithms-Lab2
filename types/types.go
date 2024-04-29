package types

type Point struct {
	X int
	Y int
}

type Rectangle struct {
	LowerLeft  Point
	UpperRight Point
}

type Event struct {
	N     int
	Left  int
	Right int
	State int
}

type Node struct {
	Left       *Node
	Right      *Node
	LeftRange  int
	RightRange int
	Sum        int
}

func NewPoint(x int, y int) *Point {
	return &Point{X: x, Y: y}
}

func NewRectangle(x1, y1, x2, y2 int) *Rectangle {
	return &Rectangle{LowerLeft: *NewPoint(x1, y1), UpperRight: *NewPoint(x2, y2)}
}

func NewEvent(n, left, right, isBegOrEnd int) *Event {
	return &Event{N: n, Left: left, Right: right, State: isBegOrEnd}
}

func NewNode(left *Node, right *Node, leftRange, rightRange, sum int) *Node {
	return &Node{Left: left, Right: right, LeftRange: leftRange,
		RightRange: rightRange, Sum: sum}
}

func (rectangle Rectangle) Belongs(point Point) bool {
	return point.X >= rectangle.LowerLeft.X && point.X <= rectangle.UpperRight.X &&
		point.Y >= rectangle.LowerLeft.Y && point.Y <= rectangle.UpperRight.Y
}
