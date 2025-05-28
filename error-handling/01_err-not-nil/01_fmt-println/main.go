package main

import (
	"fmt"
	"os"
)

func main() {
	_, err := os.Open("no-such-file.txt")
	if err != nil {
		fmt.Printf("出错了:%s\n", err)
	}
}
