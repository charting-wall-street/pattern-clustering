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

	algorithms := db.Algorithms()
	for _, alg := range algorithms {
		fmt.Printf("%s -> %s\n", alg, GetMD5Hash(alg))
	}
}
