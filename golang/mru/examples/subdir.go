package main

import (
	"fmt"
	"os"
	"path/filepath"
)

//func Walk(root string, walkFn WalkFunc) error
//type WalkFunc func(path string, info os.FileInfo, err error) error
func main() {
	err := filepath.Walk("./asset/cases/", func(path string, info os.FileInfo, err error) error {
		fmt.Println(info.Name())
		fmt.Println(filepath.Abs(info.Name()))
		return nil
	})

	if err != nil {
		panic(err)
	}
}
