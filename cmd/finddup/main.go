package main

import (
	"fmt"
	"log"
	"os"

	"github.com/vps2/finddup/internal/fs"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Missing directory argument.")
	}

	dir := os.Args[1]

	fi, err := os.Stat(dir)
	if err != nil {
		log.Fatal(err)
	}
	if !fi.IsDir() {
		log.Fatal("You didn't pass the directory in the argument.")
	}

	duplicateDetector := fs.NewDuplicateDetector(dir)
	duplicates := duplicateDetector.Search()
	for k, v := range duplicates {
		fmt.Printf("Duplicate files with hash [%s]:\n", k)
		for _, path := range v.Paths {
			fmt.Println(path)
		}
		fmt.Println()
	}
}
