package utils

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ListRequestFiles scans the "webapp/requests/" directory and returns all .txt files
func ListRequestFiles() ([]string, error) {
	var files []string
	root := "webapp/requests"

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".txt") {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}
