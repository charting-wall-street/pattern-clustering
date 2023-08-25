package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"pattern-reduce/internal/distance"
	"pattern-reduce/internal/pattern"
	"pattern-reduce/internal/store"
	"pattern-reduce/internal/tree"
	"sync"
	"sync/atomic"
	"time"
)

func main() {

	input := store.PatternsFromDisk("./output/2_sampled-patterns.gob")
	fmt.Printf("total patterns: %d\n", len(input))
	n := input[0].Length()

	rootCluster := &pattern.Cluster{
		Patterns:  make([]*pattern.Tracker, 0),
		Frequency: len(input),
	}
	for _, p := range input {
		rootCluster.Patterns = append(rootCluster.Patterns, pattern.NewTracker(p))
	}

	root := tree.NewRoot(rootCluster)
	queue := []*tree.Node{root}
	queueLock := sync.Mutex{}
	busy := atomic.Int32{}

	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		j := i
		go func() {
			for {
				d := distance.NewEditDistance(n, 5)
				queueLock.Lock()
				if len(queue) == 0 {
					if busy.Load() == 0 {
						queueLock.Unlock()
						break
					} else {
						queueLock.Unlock()
						time.Sleep(10 * time.Millisecond)
						continue
					}
				}
				next := queue[0]
				queue = queue[1:]
				remaining := len(queue)
				busy.Add(1)
				queueLock.Unlock()
				hasSplit := next.Split(d, j, remaining)
				if hasSplit {
					queueLock.Lock()
					if len(next.Left.Cluster.Patterns) > 15 {
						queue = append(queue, next.Left)
					}
					if len(next.Right.Cluster.Patterns) > 15 {
						queue = append(queue, next.Right)
					}
					queueLock.Unlock()
				}
				busy.Add(-1)
			}
			wg.Done()
		}()
	}

	wg.Wait()

	f, err := os.Create("./output/4_pattern-tree.gob")
	if err != nil {
		log.Fatalln(err)
	}
	err = gob.NewEncoder(f).Encode(root)
	if err != nil {
		log.Fatalln(err)
	}
	_ = f.Close()
}
