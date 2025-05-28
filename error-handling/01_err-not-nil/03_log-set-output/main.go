package main

import (
	"log"
	"os"
)

func init() {
	file, err := os.Create("log.txt")
	if err != nil {
		log.Println(err)
	}
	//这里类似于一个回调
	log.SetOutput(file)
}
func main() {
	_, err := os.Open("no-such-file.txt")
	if err != nil {
		log.Println("出错了:", err)
	}
}
