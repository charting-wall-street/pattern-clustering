package store

import (
	"encoding/gob"
	"log"
	"os"
	"pattern-reduce/internal/pattern"
)

type PatternDiskStore struct {
	Patterns []*pattern.Pattern
}

func PatternsToDisk(patterns []*pattern.Pattern, path string) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatalln(err)
	}
	payload := &PatternDiskStore{
		Patterns: patterns,
	}
	err = gob.NewEncoder(file).Encode(payload)
	if err != nil {
		log.Println("failed to write to disk")
		log.Fatalln(err)
	}
	_ = file.Close()
}

func PatternsFromDisk(path string) []*pattern.Pattern {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	payload := new(PatternDiskStore)
	err = gob.NewDecoder(file).Decode(payload)
	if err != nil {
		log.Println("failed to read from disk")
		log.Fatalln(err)
	}
	_ = file.Close()
	return payload.Patterns
}

type ClusterDiskStore struct {
	Clusters []*pattern.Cluster
}

func ClustersToDisk(clusters []*pattern.Cluster, path string) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatalln(err)
	}
	payload := &ClusterDiskStore{
		Clusters: clusters,
	}
	err = gob.NewEncoder(file).Encode(payload)
	if err != nil {
		log.Println("failed to write to disk")
		log.Fatalln(err)
	}
	_ = file.Close()
}

func ClustersFromDisk(path string) []*pattern.Cluster {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	payload := new(ClusterDiskStore)
	err = gob.NewDecoder(file).Decode(payload)
	if err != nil {
		log.Println("failed to read from disk")
		log.Fatalln(err)
	}
	_ = file.Close()
	return payload.Clusters
}
