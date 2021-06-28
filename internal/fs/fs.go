package fs

import (
	"context"
	"os"

	"github.com/vps2/finddup/internal/hash"
)

type Duplicates struct {
	Count int
	Paths []string
}

type DuplicateDetector interface {
	Search(context.Context) map[string]*Duplicates
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

func isContextDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
