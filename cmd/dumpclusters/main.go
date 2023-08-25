package main

import (
	"fmt"
	"log"
	"os"
	"pattern-reduce/internal/store"
	"strconv"
	"strings"
)

func main() {

	if len(os.Args) != 2 {
		log.Fatalln("error: no index argument provided")
	}
	indices := os.Args[1]

	xs := strings.Split(indices, ":")
	if len(xs) > 2 {
		log.Fatalln("invalid index")
	}

	clusters := store.ClustersFromDisk("./output/clusters.bin")
	if len(xs) == 1 {
		i, err := strconv.Atoi(xs[0])
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("Cluster %d | Frequency %d\n", i, clusters[i].Frequency)
		fmt.Println(clusters[i].Patterns[0].Pattern.String())
	} else {
		start, err := strconv.Atoi(xs[0])
		if err != nil {
			log.Fatalln(err)
		}
		end, err := strconv.Atoi(xs[1])
		if err != nil {
			log.Fatalln(err)
		}
		for i := start; i < end; i++ {
			fmt.Printf("Cluster %d | Frequency %d\n", i, clusters[i].Frequency)
			fmt.Println(clusters[i].Patterns[0].Pattern.String())
		}
	}

}
