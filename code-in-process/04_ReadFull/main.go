package main

import (
	"io"
	"os"
)

func main() {
	src, err := os.Open("src.txt")
	if err != nil {
		panic(err)
	}
	defer src.Close()

	dst3, err := os.Create("dst3.txt")
	if err != nil {
		panic(err)
	}
	defer dst3.Close()

	//bs := make([]byte, 5)
	//io.ReadFull(src, bs)
	//dst3.Write(bs)
	reader := io.LimitReader(src, 5)
	io.Copy(dst3, reader)
}
