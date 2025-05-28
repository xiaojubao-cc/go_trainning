package main

import (
	"os"
)

func main() {
	_, err := os.Open("no-such-file.txt")
	if err != nil {
		//直接终止程序
		panic(err)
	}
}
