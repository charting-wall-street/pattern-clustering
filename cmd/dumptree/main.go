package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"pattern-reduce/internal/tree"
)

func main() {

	f, err := os.Open("./output/4_pattern-tree.gob")
	if err != nil {
		log.Fatalln(err)
	}
	root := new(tree.Node)
	err = gob.NewDecoder(f).Decode(root)
	if err != nil {
		log.Fatalln(err)
	}
	_ = f.Close()

	fmt.Println(root)
}
