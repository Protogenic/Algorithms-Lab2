# Лабораторная работа №2
### Описание задачи
> Даны прямоугольники на плоскости с углами в целочисленных координатах. Требуется как можно быстрее выдавать ответ на вопрос «Скольким прямоугольникам принадлежит точка x, y?». Если точка находится на границе прямоугольника, то считается, что она принадлежит ему. Подготовка данных должна занимать как можно меньше времени.

### Описание проекта
**contest:** Код одним файлом отправленный на контест.

**solution:** Папка с решениями.
- *bruteForce* - алгоритм перебора
- compressedMap - алгоритм на карте
- segmentTree - алгоритм на дереве

**types:** Структуры данных используемые в проекте.


## 1. Алгоритм перебора
Алгоритм заключается в переборе точек и поочередной проверке принадлежности каждой точки к каждому прямоугольнику.

```ruby
func BruteForce(rectangles []types.Rectangle, points []types.Point) []int {
	result := make([]int, len(points))
	for i := 0; i < len(points); i++ {
		for j := 0; j < len(rectangles); j++ {
			if rectangles[j].Belongs(points[i]) {
				result[i]++
			}
		}
	}

	return result
}
```


| Preprocessing | Query |
|----------|----------|
| O(1)  | O(n)   |

## 2. Алгоритм на карте
Алгоритм сжимает координаты прямуугольников по осям и используя бинарный поиск находит ответ для каждой точки на построенной карте.

```ruby
func CompressedMap(rectangles []types.Rectangle, points []types.Point) []int {
	compressedX, compressedY, compressedMap := CompressedMapUtil(rectangles)

	result := make([]int, len(points))
	for i := range points {
		X := binarySearch(compressedX, points[i].X)
		Y := binarySearch(compressedY, points[i].Y)

		if X == -1 || Y == -1 {
			continue
		} else {
			result[i] = compressedMap[Y][X]
		}
	}

	return result
}

func CompressedMapUtil(rectangles []types.Rectangle) ([]int, []int, [][]int) {
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

	return compressedX, compressedY, compressedMap
}
```


| Preprocessing | Query |
|----------|----------|
| O(n^3)  | O(log(n))   |

## 3. Алгоритм на дереве
Алгоритм заключается в построении персистентного дерева для оптимизации препроцессинга второго решения.

```ruby
func SegmentTree(rectangles []types.Rectangle, points []types.Point) []int {
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

	events := make([]types.Event, 0, 2*len(rectangles))
	for _, rectangle := range rectangles {
		events = append(events, *types.NewEvent(BinarySearch(compressedX, rectangle.LowerLeft.X),
			BinarySearch(compressedY, rectangle.LowerLeft.Y),
			BinarySearch(compressedY, rectangle.UpperRight.Y)-1, 1))
		events = append(events, *types.NewEvent(BinarySearch(compressedX, rectangle.UpperRight.X),
			BinarySearch(compressedY, rectangle.LowerLeft.Y),
			BinarySearch(compressedY, rectangle.UpperRight.Y)-1, -1))
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].N < events[j].N
	})

	values := make([]int, len(compressedY))
	root := BuildTree(values, 0, len(compressedY)-1)
	roots := make([]*types.Node, 0, 2*len(rectangles)+1)
	lastX := events[0].N
	for _, event := range events {
		if event.N != lastX {
			roots = append(roots, root)
			lastX = event.N
		}
		root = AddNode(root, event.Left, event.Right, event.State)
	}

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
	return result
}

func BuildTree(array []int, leftIndex, rightIndex int) *types.Node {
	if leftIndex >= rightIndex {
		return types.NewNode(nil, nil, leftIndex, rightIndex, array[leftIndex])
	}
	middle := (leftIndex + rightIndex) / 2

	left := BuildTree(array, leftIndex, middle)
	right := BuildTree(array, middle+1, rightIndex)

	return types.NewNode(left, right, left.LeftRange, right.RightRange, left.Sum+right.Sum)
}

func AddNode(root *types.Node, leftIndex, rightIndex, sum int) *types.Node {
	if leftIndex <= root.LeftRange && root.RightRange <= rightIndex {
		return types.NewNode(root.Left, root.Right, root.LeftRange, root.RightRange, root.Sum+sum)
	}
	if root.RightRange < leftIndex || rightIndex < root.LeftRange {
		return root
	}
	node := types.NewNode(root.Left, root.Right, root.LeftRange, root.RightRange, root.Sum)
	node.Left = AddNode(node.Left, leftIndex, rightIndex, sum)
	node.Right = AddNode(node.Right, leftIndex, rightIndex, sum)
	return node
}
```


| Preprocessing | Query |
|----------|----------|
| O(log(n))  | O(log(n))   |





