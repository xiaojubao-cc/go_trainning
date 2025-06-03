package main

import (
	"bufio"
	"os"
)

func main() {
	openFile, _ := os.OpenFile("src.txt", os.O_RDWR|os.O_CREATE, 0666)
	defer openFile.Close()
	/*逐行读取数据*/
	scanner := bufio.NewScanner(openFile)
	for scanner.Scan() {
		println(scanner.Text())
	}
}
