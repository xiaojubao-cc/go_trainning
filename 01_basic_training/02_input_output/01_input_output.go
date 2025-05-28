package main

import (
	"fmt"
	"os"
)

type person struct {
	name    string
	age     int
	address string
}

func main() {
	os.Stdout.Write([]byte("hello world\n"))
	/*指针输出*/
	var str string = "hello world"
	var prt *string = &str
	fmt.Printf("输出指针：%p\n", prt)
	fmt.Printf("%v\n", person{"lihua", 22, "beijing"})
	fmt.Printf("%+v\n", person{"lihua", 22, "beijing"})
	fmt.Printf("%#v\n", person{"lihua", 22, "beijing"})
}
