package main

import (
	"log"
	"os"
)

func main() {
	_, err := os.Open("no-such-file.txt")
	if err != nil {
		log.Println("出错了:", err)
	}
}
