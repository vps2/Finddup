package fs

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/vps2/finddup/internal/hash"
)

type DuplicateDetector struct {
	Root string
}

func NewDuplicateDetector(root string) *DuplicateDetector {
	return &DuplicateDetector{
		Root: root,
	}
}

func (d *DuplicateDetector) Search() map[string]*Duplicates {
	duplicates := make(map[string]*Duplicates)

	_ = filepath.WalkDir(d.Root, func(path string, d fs.DirEntry, e error) error {
		if e != nil {
			log.Println(e)
			return nil
		}

		if d.IsDir() {
			return nil
		}

		fHash, err := calculateFileHash(path)
		if err != nil {
			log.Println(err)
			return nil
		}

		entry, ok := duplicates[fHash]
		if !ok {
			entry = &Duplicates{}
			duplicates[fHash] = entry
		}
		entry.Count++
		entry.Paths = append(entry.Paths, path)

		return nil
	})

	for k, v := range duplicates {
		if v.Count < 2 {
			delete(duplicates, k)
		}
	}

	return duplicates
}

func calculateFileHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fHash, err := hash.Calculate(file)
	if err != nil {
		return "", err
	}

	return fHash, nil
}
