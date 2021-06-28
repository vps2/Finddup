package fs

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type fileInfo struct {
	path string
	hash string
}

type DuplicateDetectorParallel struct {
	Root string
}

func NewDuplicateDetectorParallel(root string) *DuplicateDetectorParallel {
	return &DuplicateDetectorParallel{
		Root: root,
	}
}

//FIXME сейчас запускается неограниченное количество горутин. На обычных дисках это создает проблему с зависанием.
func (d *DuplicateDetectorParallel) Search(ctx context.Context) map[string]*Duplicates {
	resultsCh := make(chan fileInfo)
	errorsCh := make(chan error)

	var wg sync.WaitGroup
	wg.Add(1)
	go visit(ctx, d.Root, &wg, resultsCh, errorsCh)

	go func() {
		wg.Wait()
		close(errorsCh)
		close(resultsCh)
	}()

	go func() {
		for error := range errorsCh {
			log.Println(error)
		}
	}()

	duplicates := make(map[string]*Duplicates)
	for fi := range resultsCh {
		entry, ok := duplicates[fi.hash]
		if !ok {
			entry = &Duplicates{}
			duplicates[fi.hash] = entry
		}
		entry.Count++
		entry.Paths = append(entry.Paths, fi.path)

		if isContextDone(ctx) {
			break
		}
	}

	for k, v := range duplicates {
		if v.Count < 2 {
			delete(duplicates, k)
		}
	}

	return duplicates
}

func visit(ctx context.Context, root string, wg *sync.WaitGroup, results chan<- fileInfo, errors chan<- error) {
	defer wg.Done()

	if isContextDone(ctx) {
		return
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		errors <- err
	}

	for _, entry := range entries {
		entryPath := filepath.Join(root, entry.Name())
		switch {
		case entry.IsDir():
			wg.Add(1)
			go visit(ctx, entryPath, wg, results, errors)
		case entry.Type().IsRegular():
			hash, err := calculateFileHash(entryPath)
			if err != nil {
				errors <- err
			}
			results <- fileInfo{path: entryPath, hash: hash}
		}

		if isContextDone(ctx) {
			return
		}
	}
}
