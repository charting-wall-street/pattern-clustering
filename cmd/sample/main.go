package main

import (
	"fmt"
	"math/rand"
	"pattern-reduce/internal/pattern"
	"pattern-reduce/internal/store"
	"strconv"
)

func RandomSample(slice []*pattern.Pattern, size int) []*pattern.Pattern {
	length := len(slice)
	if size > length {
		size = length
	}
	copiedSlice := make([]*pattern.Pattern, length)
	copy(copiedSlice, slice)
	for i := length - 1; i >= 1; i-- {
		j := rand.Intn(i + 1)
		copiedSlice[i], copiedSlice[j] = copiedSlice[j], copiedSlice[i]
	}
	return copiedSlice[:size]
}

func main() {
	inPatterns := store.PatternsFromDisk("./output/1_split-patterns.gob")
	fmt.Printf("total patterns: %d\n", len(inPatterns))
	outPatterns := RandomSample(inPatterns, len(inPatterns)/10)
	fmt.Printf("sampled patterns: %d\n", len(outPatterns))

	// sanity check
	for _, p := range outPatterns {
		if len(p.Sequences) != 3 {
			panic("fail: got " + strconv.Itoa(len(p.Sequences)))
		}
		for _, sequence := range p.Sequences {
			if sequence.Events == nil {
				panic("fail")
			}
		}
	}

	store.PatternsToDisk(outPatterns, "./output/2_sampled-patterns.gob")
}
