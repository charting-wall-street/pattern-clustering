package distance

import (
	"fmt"
	"testing"
)

func TestSimpleDistance(t *testing.T) {

	d := NewEditDistance(5, 2)

	dist := d.Distance([]float32{0, 1, 0, 0, 0}, []float32{0, 0, 0, 1, 0})

	fmt.Println(dist)
}
