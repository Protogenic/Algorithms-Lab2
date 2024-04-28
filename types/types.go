package types

type Point struct {
	X int
	Y int
}

func NewPoint(x int, y int) *Point {
	return &Point{X: x, Y: y}
}

type Rectangle struct {
	LowerLeft  Point
	UpperRight Point
}

func NewRectangle(x1, y1, x2, y2 int) *Rectangle {
	return &Rectangle{LowerLeft: *NewPoint(x1, y1), UpperRight: *NewPoint(x2, y2)}
}

func (rectangle Rectangle) Belongs(point Point) bool {
	return point.X >= rectangle.LowerLeft.X && point.X <= rectangle.UpperRight.X &&
		point.Y >= rectangle.LowerLeft.Y && point.Y <= rectangle.UpperRight.Y
}

type Event struct {
	N          int
	Left       int
	Right      int
	IsBegOrEnd int
}

type Node struct {
	Left       *Node
	Right      *Node
	LeftRange  int
	RightRange int
	Sum        int
}

func NewNode(left *Node, right *Node, leftRange, rightRange, sum int) *Node {
	return &Node{Left: left, Right: right, LeftRange: leftRange,
		RightRange: rightRange, Sum: sum}
}

type PersistentSegmentTree struct {
	Nodes []Node
}

func NewPersistentSegmentTree(events []Event, versions []*Node, size int) []*Node {
	arr := make([]int, size)
	tree := buildTree(arr, 0, size)
	n := events[0].N

	for _, event := range events {
		if n != event.N {
			versions = append(versions, tree)
			n = event.N
		}
		tree = addNode(tree, event.Left, event.Right, event.IsBegOrEnd)
	}

	return versions
}

func buildTree(arr []int, leftIndex, rightIndex int) *Node {
	if rightIndex-leftIndex == 1 {
		return NewNode(nil, nil, leftIndex, rightIndex, arr[leftIndex])
	}

	middle := (leftIndex + rightIndex) / 2

	left := buildTree(arr, leftIndex, middle)
	right := buildTree(arr, middle+1, rightIndex)

	return NewNode(left, right, left.LeftRange, right.RightRange, left.Sum+right.Sum)
}

func addNode(root *Node, leftIndex, rightIndex, val int) *Node {
	if leftIndex <= root.LeftRange && rightIndex >= root.RightRange {
		return NewNode(root.Left, root.Right, root.LeftRange, root.RightRange, root.Sum+val)
	}

	if root.LeftRange >= rightIndex || root.RightRange <= leftIndex {
		return root
	}

	newRoot := NewNode(root.Left, root.Right, root.LeftRange, root.RightRange, root.Sum)

	newRoot.Left = addNode(newRoot.Left, leftIndex, rightIndex, val)
	newRoot.Right = addNode(newRoot.Right, leftIndex, rightIndex, val)

	return newRoot
}
