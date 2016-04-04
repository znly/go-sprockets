package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/znly/go-sprockets"
)

func main() {
	s, err := sprockets.NewWithDefault(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Println(err)
		return
	}
	filepath.Walk(os.Args[1], func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if ext == ".coffee" {
			path = path[0:len(path)-len(ext)] + ".js"
		} else if ext == ".sass" || ext == ".scss" {
			path = path[0:len(path)-len(ext)] + ".css"
		}
		relPath, err := filepath.Rel(os.Args[1], path)
		if err != nil {
			fmt.Printf("Cant Generate %s: %s\n", path, err)
			return nil
		}
		if err := s.Generate(relPath); err != nil {
			fmt.Println(err)
		}
		return nil
	})
}
