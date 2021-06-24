package hash

import (
	"crypto/sha1"
	"fmt"
	"io"
)

func Calculate(reader io.Reader) (string, error) {
	hash := sha1.New()
	_, err := io.Copy(hash, reader)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
