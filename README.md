# Лабораторная работа №2
### Описание задачи
> Даны прямоугольники на плоскости с углами в целочисленных координатах. Требуется как можно быстрее выдавать ответ на вопрос «Скольким прямоугольникам принадлежит точка x, y?». Если точка находится на границе прямоугольника, то считается, что она принадлежит ему. Подготовка данных должна занимать как можно меньше времени.

### Описание проекта
**contest:** Код одним файлом отправленный на контест.

**solution:** Папка с решениями.
- *bruteForce* - алгоритм перебора
- *compressedMap* - алгоритм на карте
- *segmentTree* - алгоритм на дереве

**types:** Структуры данных используемые в проекте.

**raw:** Данные, использованные в графиках.


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

## Тестирование
Код тестировался на заданном количестве (n) случайно сгенерированных прямоугольников и точек. 0 <= n 2^20
```ruby
func randomData() ([]types.Rectangle, []types.Point) {
	var numberOfRectangles, numberOfPoints int
	fmt.Scanf("%d %d", &numberOfRectangles, &numberOfPoints)
	max := int(1000000000)

	rectangles := make([]types.Rectangle, numberOfRectangles)
	for i := 0; i < numberOfRectangles; i++ {
		var x1, y1, x2, y2 int
		x1 = rand.Intn(max)
		x2 = rand.Intn(max-x1) + x1 + 1
		y1 = rand.Intn(max)
		y2 = rand.Intn(max-y1) + y1 + 1
		rectangles[i] = *types.NewRectangle(x1, y1, x2, y2)
	}

	points := make([]types.Point, numberOfPoints)
	for i := 0; i < numberOfPoints; i++ {
		var x, y int
		x = rand.Intn(max)
		y = rand.Intn(max)
		points[i] = *types.NewPoint(x, y)
	}

	return rectangles, points
}
```
---
## Графики
### Общее время выполнения
![image](https://github.com/Protogenic/Algorithms-Lab2/assets/82569672/836258db-81f2-40ed-abc1-a7d1d2fbf36d)

![image](https://github.com/Protogenic/Algorithms-Lab2/assets/82569672/7e549a5d-6516-4122-98c5-dfb2bc957cb5)

Из этих графиков видно, что время выполнения первого алгоритма, из-за большего времени затраченного на выполнение одного запроса, растет гораздо быстрее с увеличением количества точек, чем в случае применения двух других алгоритмов.

### Препроцессинг и запросы
![image](https://github.com/Protogenic/Algorithms-Lab2/assets/82569672/7dd07f2a-0e40-4f60-b661-9df9bd6a3b0b)

![image](https://github.com/Protogenic/Algorithms-Lab2/assets/82569672/8edc648c-76c8-436e-a5a4-882b4ec2be74)

Из приведенных выше графиков можно заметить, что время потраченное на обработку запросов у 2 и 3 алгоритма схожи, но алгоритм на карте сильно уступает в общем времени из-за крайне долгого препроцессинга.

## Выводы
Алгоритм перебора подходит для небольших данных, особенно в случае малого количества точек, из-за отсутствия времени затраченного на препроцессинг на таких входных данных этот алгоритм будет выигрывать по времени у двух других.

Алгоритм на карте будет хорошо работать в случае малого количества прямоугольников и большого количества точек. Препроцессинг работает за O(n^3), он повлечет огромные затраты по времени в случае большого количества прямоугольников, а выполнение запроса, в отличие от первого алгоритма, выполняется за O(log(n)), поэтому увеличение количества точек на этот алгоритм влияет незначительно.

Алгоритм на дереве отлично подходит для больших входных данных, так как препроцессинг не будет занимать колоссальное количество времени, как в случае алгоритма на карте, и каждый запрос будет быстро выполняться, в отличии от алгоритма перебора.


