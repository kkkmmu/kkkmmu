package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	var files []string
	filepath.Walk("asset/cases", func(path string, info os.FileInfo, err error) error {
		if info.Name() == "all.yaml" {
			files = append(files, path)
		}
		return nil
	})

	for _, file := range files {
		fmt.Println(file)
	}
}
