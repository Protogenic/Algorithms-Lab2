package contest

import (
	"fmt"
	"sort"
)

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

func SegmentTree(rectangles []Rectangle, points []Point) []int {
	//fmt.Print("Preparation time: ")
	//startTime := time.Now()

	if len(rectangles) == 0 {
		return make([]int, 0)
	}

	uniqueX := make(map[int]struct{})
	uniqueY := make(map[int]struct{})
	for _, rectangle := range rectangles {
		uniqueX[rectangle.LowerLeft.X] = struct{}{}
		uniqueX[rectangle.UpperRight.X] = struct{}{}

		uniqueY[rectangle.LowerLeft.Y] = struct{}{}
		uniqueY[rectangle.UpperRight.Y] = struct{}{}
	}

	var compressedX []int
	for x := range uniqueX {
		compressedX = append(compressedX, x)
	}
	sort.Slice(compressedX, func(i, j int) bool {
		return compressedX[i] < compressedX[j]
	})

	var compressedY []int
	for y := range uniqueY {
		compressedY = append(compressedY, y)
	}
	sort.Slice(compressedY, func(i, j int) bool {
		return compressedY[i] < compressedY[j]
	})

	events := make([]Event, 0, 2*len(rectangles))
	for _, rectangle := range rectangles {
		events = append(events, *NewEvent(BinarySearch(compressedX, rectangle.LowerLeft.X),
			BinarySearch(compressedY, rectangle.LowerLeft.Y),
			BinarySearch(compressedY, rectangle.UpperRight.Y)-1, 1))
		events = append(events, *NewEvent(BinarySearch(compressedX, rectangle.UpperRight.X),
			BinarySearch(compressedY, rectangle.LowerLeft.Y),
			BinarySearch(compressedY, rectangle.UpperRight.Y)-1, -1))
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].N < events[j].N
	})

	values := make([]int, len(compressedY))
	root := BuildTree(values, 0, len(compressedY)-1)
	roots := make([]*Node, 0, 2*len(rectangles)+1)
	lastX := events[0].N
	for _, event := range events {
		if event.N != lastX {
			roots = append(roots, root)
			lastX = event.N
		}
		root = AddNode(root, event.Left, event.Right, event.State)
	}

	//endTime := time.Now()
	//duration := endTime.Sub(startTime)
	//fmt.Println(duration)

	//startTime = time.Now()

	result := make([]int, len(points))
	index := 0
	for _, point := range points {
		xPos := BinarySearch(compressedX, point.X)
		yPos := BinarySearch(compressedY, point.Y)
		if xPos == -1 || yPos == -1 || xPos >= len(roots) {
			index++
			continue
		}
		result[index] = BinarySearchTree(roots[xPos], yPos)
		index++
	}

	//endTime = time.Now()
	//duration = endTime.Sub(startTime)
	//fmt.Println("Execution time: ", duration, "\n")

	for i := 0; i < len(result); i++ {
		fmt.Print(result[i])
		if i != len(result)-1 {
			fmt.Print(" ")
		}
	}

	return result
}

func BuildTree(array []int, leftIndex, rightIndex int) *Node {
	if leftIndex >= rightIndex {
		return NewNode(nil, nil, leftIndex, rightIndex, array[leftIndex])
	}
	middle := (leftIndex + rightIndex) / 2

	left := BuildTree(array, leftIndex, middle)
	right := BuildTree(array, middle+1, rightIndex)

	return NewNode(left, right, left.LeftRange, right.RightRange, left.Sum+right.Sum)
}

func AddNode(root *Node, leftIndex, rightIndex, sum int) *Node {
	if leftIndex <= root.LeftRange && root.RightRange <= rightIndex {
		return NewNode(root.Left, root.Right, root.LeftRange, root.RightRange, root.Sum+sum)
	}
	if root.RightRange < leftIndex || rightIndex < root.LeftRange {
		return root
	}
	node := NewNode(root.Left, root.Right, root.LeftRange, root.RightRange, root.Sum)
	node.Left = AddNode(node.Left, leftIndex, rightIndex, sum)
	node.Right = AddNode(node.Right, leftIndex, rightIndex, sum)
	return node
}

func BinarySearch(array []int, target int) int {
	left, right := 0, len(array)-1
	for left <= right {
		middle := (left + right) / 2
		if array[middle] == target {
			return middle
		}
		if target < array[middle] {
			right = middle - 1
		} else {
			left = middle + 1
		}
	}

	return right
}

func BinarySearchTree(root *Node, index int) int {
	if root == nil {
		return 0
	}

	middle := (root.LeftRange + root.RightRange) / 2

	var sum int
	if index <= middle {
		sum = BinarySearchTree(root.Left, index)
	} else {
		sum = BinarySearchTree(root.Right, index)
	}
	return sum + root.Sum
}

func readData() ([]Rectangle, []Point) {
	var numberOfRectangles int
	fmt.Scanf("%d\n", &numberOfRectangles)
	rectangles := make([]Rectangle, numberOfRectangles)
	for i := 0; i < numberOfRectangles; i++ {
		var x1, y1, x2, y2 int
		fmt.Scanf("%d %d %d %d\n", &x1, &y1, &x2, &y2)
		rectangles[i] = *NewRectangle(x1, y1, x2, y2)
	}

	var numberOfPoints int
	fmt.Scanf("%d\n", &numberOfPoints)
	points := make([]Point, numberOfPoints)
	for i := 0; i < numberOfPoints; i++ {
		var x, y int
		fmt.Scanf("%d %d\n", &x, &y)
		points[i] = *NewPoint(x, y)
	}

	return rectangles, points
}

func main() {
	rectangles, points := readData()
	SegmentTree(rectangles, points)
}
