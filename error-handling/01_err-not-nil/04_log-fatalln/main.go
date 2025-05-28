package main

import (
	"log"
	"os"
)

func main() {
	_, err := os.Open("no-such-file.txt")
	if err != nil {
		//打印日志信息并终止程序
		log.Fatalln("出错了:", err)
	}
}
