package main

import (
	"fmt"
	"os"

	sprockets "github.com/znly/go-sprockets"
)

func main() {
	s, err := sprockets.NewWithDefault(os.Args[1], "")
	if err != nil {
		fmt.Println(err)
		return
	}
	content, err := s.GetAsset(os.Args[2])
	fmt.Println(string(content), err)
}
