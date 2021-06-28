package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/vps2/finddup/internal/fs"
)

func main() {
	dir := flag.String("dir", "", "The directory in which to search for duplicates.")
	parallel := flag.Bool("parallel", false, "Run in parallel. Use only with SSD.")
	flag.Parse()

	var duplicateDetector fs.DuplicateDetector
	if *dir == "" {
		flag.Usage()
		return
	}
	if *parallel {
		duplicateDetector = fs.NewDuplicateDetectorParallel(*dir)
	} else {
		duplicateDetector = fs.NewDuplicateDetectorSerial(*dir)
	}

	fi, err := os.Stat(*dir)
	if err != nil {
		log.Fatal(err)
	}
	if !fi.IsDir() {
		log.Fatal("You didn't pass the directory in the argument.")
	}

	ctx, cancel := context.WithCancel(context.Background())

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt)

	go func() {
		<-stopCh
		cancel()
	}()

	duplicates := duplicateDetector.Search(ctx)
	for k, v := range duplicates {
		fmt.Printf("Duplicate files with hash [%s]:\n", k)
		for _, path := range v.Paths {
			fmt.Println(path)
		}
		fmt.Println()
	}
}
