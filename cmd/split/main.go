package main

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/godoji/algocore/pkg/algo"
	"github.com/godoji/algocore/pkg/kiosk"
	"github.com/northberg/candlestick"
	"log"
	"math"
	"pattern-reduce/internal/db"
	"pattern-reduce/internal/pattern"
	"pattern-reduce/internal/store"
	"sync"
)

func fetchEvents(algorithm string, symbol string) *algo.ScenarioSet {
	result, err := kiosk.GetAlgorithm(algorithm, candlestick.Interval1d, symbol, []float64{}, true)
	if err != nil {
		panic(err)
	}
	return result
}

func createTimeline(firstEvent int64, lastEvent int64, events map[string]*algo.ScenarioSet) []*pattern.Pattern {
	lastIndex := make(map[string]int)
	slotCount := (lastEvent-firstEvent)/candlestick.Interval1d + 1
	slots := make(map[string][]float32)
	for key := range events {
		slots[key] = make([]float32, slotCount)
	}
	for i := firstEvent; i <= lastEvent; i += candlestick.Interval1d {
		for algorithm, set := range events {
			for index, event := range set.Events[lastIndex[algorithm]:] {
				if event.CreatedOn == i {
					lastIndex[algorithm] = index
					slotIndex := (i - firstEvent) / candlestick.Interval1d
					if event.Label == "uptrend" {
						slots[algorithm][slotIndex] = 1
					} else {
						slots[algorithm][slotIndex] = -1
					}
				}
			}
		}
	}

	windowSize := int64(14)
	windowOverlap := int64(2)
	patterns := make([]*pattern.Pattern, 0)
	for i := int64(0); i < slotCount/windowSize-windowOverlap; i++ {
		sequences := make(map[string]pattern.EventSequence)
		for key, evs := range slots {
			xs := evs[windowSize*i : windowSize*(i+windowOverlap)]
			empty := true
			for j := 0; j < len(xs); j++ {
				if xs[j] != 0 {
					empty = false
					break
				}
			}

			if len(xs) == 0 || xs == nil {
				log.Fatalln("event list should not be empty or nil")
			}

			sequences[key] = pattern.EventSequence{
				Empty:  empty,
				Events: xs,
			}
		}
		collection := &pattern.Pattern{
			Sequences: sequences,
		}
		patterns = append(patterns, collection)
	}

	return patterns
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func main() {

	algorithms := []string{"ema-cross", "rsi-levels", "bollinger-bands"}
	symbols := db.Symbols()

	patterns := make([]*pattern.Pattern, 0)
	patternsLock := sync.Mutex{}

	wg := &sync.WaitGroup{}
	for _, symbol := range symbols {
		events := make(map[string]*algo.ScenarioSet)
		firstEvent := int64(math.MaxInt64)
		lastEvent := int64(math.MinInt64)
		for _, algorithm := range algorithms {
			scenario := fetchEvents(algorithm, symbol)
			if len(scenario.Events) == 0 {
				continue
			}
			events[GetMD5Hash(algorithm)] = scenario
			if firstEvent > scenario.Events[0].CreatedOn {
				firstEvent = scenario.Events[0].CreatedOn
			}
			if lastEvent < scenario.Events[len(scenario.Events)-1].CreatedOn {
				lastEvent = scenario.Events[len(scenario.Events)-1].CreatedOn
			}
		}

		if firstEvent >= lastEvent {
			continue
		}

		wg.Add(1)
		go func() {
			result := createTimeline(firstEvent, lastEvent, events)
			patternsLock.Lock()
			patterns = append(patterns, result...)
			patternsLock.Unlock()
			wg.Done()
		}()
	}

	wg.Wait()

	store.PatternsToDisk(patterns, "./output/1_split-patterns.gob")
}
