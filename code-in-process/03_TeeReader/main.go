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

	dst1, err := os.Create("dst1.txt")
	if err != nil {
		panic(err)
	}
	defer dst1.Close()

	dst2, err := os.Create("dst2.txt")
	if err != nil {
		panic(err)
	}
	defer dst2.Close()

	//创建一个TeeReader,读取文件的同时写入内容到新目录中，并且在控制台打印
	reader1 := io.TeeReader(src, dst1)
	reader2 := io.TeeReader(reader1, os.Stdout)

	io.Copy(dst2, reader2)
}
