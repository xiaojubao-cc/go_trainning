package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	filepath.Walk("D:\\golang projects\\go_training\\code-in-process\\05_review-exercises", func(path string, info os.FileInfo, err error) error {
		fmt.Println(path)
		fmt.Println(info.IsDir())
		fmt.Println(info.Name())
		fmt.Println(info.Size())
		fmt.Println(info.Mode())
		fmt.Println(info.ModTime())
		fmt.Println(info.Sys())
		fmt.Println(err)
		return nil
	})
}
