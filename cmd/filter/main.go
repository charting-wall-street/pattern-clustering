package main

import (
	"log"
	"pattern-reduce/internal/distance"
	"pattern-reduce/internal/pattern"
	"pattern-reduce/internal/store"
	"sync"
)

func ReduceEquivalentPatterns(patterns []*pattern.Pattern, d *distance.EditDistanceCalculator, tid int) []*pattern.Tracker {
	output := make([]*pattern.Tracker, 0)
	patternMask := make(map[int]bool)
	removed := 0

	for i, p := range patterns {
		if i%500 == 0 {
			log.Printf("[thread %d]\tpattern: %d | remaining: %d\n", tid, i, len(patterns)-removed)
		}
		if patternMask[i] {
			continue
		}
		patternMask[i] = true
		freq := 1
		for j := i + 1; j < len(patterns); j++ {
			if patternMask[j] {
				continue
			}
			o := patterns[j]
			if distance.IsZeroDistance(p, o, d) {
				freq++
				patternMask[j] = true
			}
		}
		output = append(output, &pattern.Tracker{
			Pattern:   p,
			Frequency: freq,
		})
		removed += freq
	}

	return output
}

func ReduceEquivalentContainers(patterns []*pattern.Tracker, d *distance.EditDistanceCalculator, tid int) []*pattern.Tracker {
	output := make([]*pattern.Tracker, 0)
	patternMask := make(map[int]bool)
	removed := 0

	for i, p := range patterns {
		if patternMask[i] {
			continue
		}
		freq := 0
		patternMask[i] = true
		removed++
		for j := i + 1; j < len(patterns); j++ {
			if patternMask[j] {
				continue
			}
			o := patterns[j]
			if distance.IsZeroDistance(p.Pattern, o.Pattern, d) {
				freq += o.Frequency
				patternMask[j] = true
				removed++
			}
		}
		p.Frequency += freq
		output = append(output, p)
		if i%500 == 0 {
			log.Printf("[thread %d] pattern: %d | remaining: %d\n", tid, i, len(patterns)-removed)
		}
	}

	return output
}

func HyperReduceEquivalentPatterns(input []*pattern.Pattern, n int, threads int) []*pattern.Tracker {
	output := make([]*pattern.Tracker, 0)
	wg := sync.WaitGroup{}
	patternsLock := sync.Mutex{}
	blockSize := len(input) / threads
	for i := 0; i < threads; i++ {
		j := i
		sub := input[i*blockSize : (i+1)*blockSize]
		if i == threads-1 {
			sub = input[i*blockSize:]
		}
		wg.Add(1)
		go func() {
			d := distance.NewEditDistance(n, 6)
			result := ReduceEquivalentPatterns(sub, d, j)
			patternsLock.Lock()
			output = append(output, result...)
			patternsLock.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
	return output
}

func HyperReduceEquivalentContainers(patterns []*pattern.Tracker, n int, threads int) []*pattern.Tracker {
	output := make([]*pattern.Tracker, 0)
	wg := sync.WaitGroup{}
	patternsLock := sync.Mutex{}
	blockSize := len(patterns) / threads
	for i := 0; i < threads; i++ {
		j := i
		sub := patterns[i*blockSize : (i+1)*blockSize]
		if i == threads-1 {
			sub = patterns[i*blockSize:]
		}
		wg.Add(1)
		go func() {
			d := distance.NewEditDistance(n, 6)
			result := ReduceEquivalentContainers(sub, d, j)
			patternsLock.Lock()
			output = append(output, result...)
			patternsLock.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
	return output
}

func main() {
	input := store.PatternsFromDisk("./output/2_sampled-patterns.gob")
	n := input[0].Length()

	nextPatterns := HyperReduceEquivalentPatterns(input, n, 12)

	patterns := HyperReduceEquivalentContainers(nextPatterns, n, 6)
	nextPatterns = patterns
	patterns = HyperReduceEquivalentContainers(nextPatterns, n, 3)
	nextPatterns = patterns
	patterns = HyperReduceEquivalentContainers(nextPatterns, n, 1)

	log.Printf("found %d unique patterns\n", len(patterns))

	clusters := make([]*pattern.Cluster, len(patterns))
	for i, t := range patterns {
		clusters[i] = pattern.NewCluster([]*pattern.Tracker{t})
	}
	store.ClustersToDisk(clusters, "./output/3_input-clusters.gob")
}
