package main

import (
	"io"
	"os"
)

func main() {
	file1, err := os.Open("D:\\golang projects\\go_training\\dst.txt")
	if err != nil {
		panic(err)
	}
	defer file1.Close()
	file2, err := os.Open("D:\\golang projects\\go_training\\dst1.txt")
	if err != nil {
		panic(err)
	}
	defer file2.Close()

	reads := make([]io.Reader, 2)
	reads[0] = file1
	reads[1] = file2
	reader := io.MultiReader(reads...)
	io.Copy(os.Stdout, reader)
}
