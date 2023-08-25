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

	patterns := store.PatternsFromDisk("./output/all.patterns.bin")
	if len(xs) == 1 {
		i, err := strconv.Atoi(xs[0])
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("Pattern %d\n", i)
		fmt.Println(patterns[i].String())
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
			fmt.Printf("Pattern %d\n", i)
			fmt.Println(patterns[i].String())
		}
	}

}
