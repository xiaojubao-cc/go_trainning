package main

import (
	"io"
	"log"
	"os"
)

/*文件操作*/
func main() {
	openFile, err := os.OpenFile("src.txt", os.O_RDWR|os.O_CREATE, 0666)
	defer openFile.Close()
	if os.IsNotExist(err) {
		log.Fatal(err)
	} else if err != nil {
		log.Fatal(err)
	} else {
		bytes, err := io.ReadAll(openFile)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(bytes))
	}
	bytes, _ := openFile.WriteString("hello golang\n")
	log.Print(bytes)

}
