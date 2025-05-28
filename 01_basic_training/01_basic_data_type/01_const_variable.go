package main

import "fmt"

/*
包：共享包中的所有的可以使用(公有的即大写字母开头的)的变量、常量和函数,匿名导包是为了执行包中的init函数
常量定义：const定义的常量需要赋初值
*/
const (
	PI = 3.14
	E  = 2.71828
)

func main() {
	var str string = "hello world"
	fmt.Printf("打印字符串变量：%s\n", str)
	fmt.Printf("打印浮点型常量：%f\n", PI)
	fmt.Printf("打印浮点型常量：%f\n", E)
}
