package main

import (
	"pattern-reduce/internal/distance"
	"pattern-reduce/internal/store"
	"testing"
)

func BenchmarkWindowSize3(b *testing.B) {
	patterns, n := store.LoadPatternGroups("./output/patterns_test.txt")
	d := distance.NewEditDistance(n, 3)
	ReduceEquivalentPatterns(patterns, d)
}

func BenchmarkWindowSize5(b *testing.B) {
	patterns, n := store.LoadPatternGroups("./output/patterns_test.txt")
	d := distance.NewEditDistance(n, 5)
	ReduceEquivalentPatterns(patterns, d)
}

func BenchmarkWindowSize8(b *testing.B) {
	patterns, n := store.LoadPatternGroups("./output/patterns_test.txt")
	d := distance.NewEditDistance(n, 8)
	ReduceEquivalentPatterns(patterns, d)
}
