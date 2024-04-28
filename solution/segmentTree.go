package solution

import (
	"Algorithms_Lab2/types"
	"fmt"
	"sort"
)

func NewEvent(n, left, right, isBegOrEnd int) *types.Event {
	return &types.Event{N: n, Left: left, Right: right, IsBegOrEnd: isBegOrEnd}
}

func SegmentTree(rectangles []types.Rectangle, points []types.Point) {
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

	events := make([]types.Event, len(rectangles)*2)
	versions := make([]*types.Node, 0, len(rectangles)*2)

	for _, rectangle := range rectangles {
		events = append(events, *NewEvent(binarySearch(compressedX, rectangle.LowerLeft.X),
			binarySearch(compressedY, rectangle.LowerLeft.Y),
			binarySearch(compressedY, rectangle.UpperRight.Y+1), 1))
		events = append(events, *NewEvent(binarySearch(compressedX, rectangle.UpperRight.X+1),
			binarySearch(compressedY, rectangle.LowerLeft.Y),
			binarySearch(compressedY, rectangle.UpperRight.Y+1), -1))
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].N < events[j].N
	})
	nodes := types.NewPersistentSegmentTree(events, versions, len(compressedY))

	for i := 0; i < len(points); i++ {
		X := binarySearch(compressedX, points[i].X)
		Y := binarySearch(compressedY, points[i].Y)

		if X == -1 || Y == -1 || len(nodes) <= X {
			fmt.Println(points[i], 0)
		} else {
			fmt.Println(points[i], searchInTree(nodes[X], Y))
		}
	}
}

func searchInTree(root *types.Node, num int) int {
	if root == nil {
		return 0
	}
	middle := (root.LeftRange + root.RightRange) / 2
	if num < middle {
		return root.Sum + searchInTree(root.Left, num)
	}
	return root.Sum + searchInTree(root.Right, num)
}
