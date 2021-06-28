package fs

import (
	"context"
	"io/fs"
	"log"
	"path/filepath"
)

type DuplicateDetectorSerial struct {
	Root string
}

func NewDuplicateDetectorSerial(root string) *DuplicateDetectorSerial {
	return &DuplicateDetectorSerial{
		Root: root,
	}
}

func (d *DuplicateDetectorSerial) Search(ctx context.Context) map[string]*Duplicates {
	duplicates := make(map[string]*Duplicates)

	_ = filepath.WalkDir(d.Root, func(path string, d fs.DirEntry, e error) error {
		if isContextDone(ctx) {
			return ctx.Err()
		}

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
