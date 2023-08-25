package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"pattern-reduce/internal/db"
)

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func main() {
	algos := db.Algorithms()

	for _, algo := range algos {
		fmt.Printf("%s -> %s\n", algo, GetMD5Hash(algo))
	}
}
