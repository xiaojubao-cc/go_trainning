package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	source, err := os.Open("src.txt")
	if err != nil {
		panic(err)
	}
	defer source.Close()

	dest, err := os.Create("dst.txt")
	if err != nil {
		panic(err)
	}
	defer dest.Close()

	//读取文件
	reader := bufio.NewReader(source)
	_, err = io.Copy(dest, reader)
	if err != nil {
		panic(err)
	}

	//逐行读取文件
	scanner := bufio.NewScanner(source)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}
}
