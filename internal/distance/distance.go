package distance

import (
	"math"
	"pattern-reduce/internal/pattern"
)

var (
	StatSkipEmpty = 0
)

func editCost(a float32, b float32) float32 {
	if a != b {
		return 1
	} else {
		return 0
	}
}

type EditDistanceCalculator struct {
	firstRow  []float32
	secondRow []float32
	size      int
	window    int
}

func NewEditDistance(size int, window int) *EditDistanceCalculator {
	return &EditDistanceCalculator{
		firstRow:  make([]float32, size+1),
		secondRow: make([]float32, size+1),
		size:      size + 1,
		window:    window,
	}
}

func minFloat3(a float32, b float32, c float32) float32 {
	if a < b {
		return minFloat2(a, c)
	} else {
		return minFloat2(b, c)
	}
}

func minFloat2(a float32, b float32) float32 {
	if a < b {
		return a
	} else {
		return b
	}
}

func min2(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func max2(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func (e *EditDistanceCalculator) init() {
	for i := 0; i < e.size; i++ {
		e.secondRow[i] = math.MaxInt
		e.firstRow[i] = math.MaxInt
	}
	e.secondRow[0] = 0
}

func (e *EditDistanceCalculator) match(x int) float32 {
	return e.firstRow[x-1]
}
func (e *EditDistanceCalculator) insertion(x int) float32 {
	return e.firstRow[x]
}

func (e *EditDistanceCalculator) deletion(x int) float32 {
	return e.secondRow[x-1]
}

func (e *EditDistanceCalculator) calc(a float32, b float32, x int) float32 {
	cost := editCost(a, b)
	return cost + minFloat3(e.match(x), e.insertion(x), e.deletion(x))
}

func (e *EditDistanceCalculator) swap() {
	e.firstRow, e.secondRow = e.secondRow, e.firstRow
}

func (e *EditDistanceCalculator) Distance(a []float32, b []float32) float32 {
	e.init()
	//fmt.Println(e.firstRow)
	for y := 1; y < e.size; y++ {
		//fmt.Println(e.secondRow)
		e.swap()
		lower := max2(1, y-e.window)
		upper := min2(e.size, y+e.window+1)
		e.secondRow[lower-1] = math.MaxInt
		for x := lower; x < upper; x++ {
			e.secondRow[x] = e.calc(a[y-1], b[x-1], x)
		}
	}
	//fmt.Println(e.secondRow)
	return e.secondRow[e.size-1]
}

func PatternDistance(a *pattern.Pattern, b *pattern.Pattern, d *EditDistanceCalculator) float32 {
	sum := float32(0)
	for i := range a.Sequences {
		seqA := a.Sequences[i]
		seqB := b.Sequences[i]

		// Skip empty patterns
		if seqA.Empty && seqB.Empty {
			StatSkipEmpty++
			continue
		}
		dist := d.Distance(seqA.Events, seqB.Events)
		sum += dist
	}
	return sum
}

func IsZeroDistance(a *pattern.Pattern, b *pattern.Pattern, d *EditDistanceCalculator) bool {
	sum := float32(0)
	for i := range a.Sequences {
		seqA := a.Sequences[i]
		seqB := b.Sequences[i]
		if seqA.Empty != seqB.Empty {
			return false
		}
		dist := d.Distance(seqA.Events, seqB.Events)
		sum += dist
	}
	return sum == 0
}
