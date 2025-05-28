package main

import (
	"fmt"
	"strings"
)

/*基本数据类型赋值：int、float、bool、string、数组、结构体属于值复制*/
func main() {
	/*本质上是一个字节数组*/
	var str string = "hello world"
	for _, value := range str {
		fmt.Printf("value=%c\n", value)
	}
	str1 := "石头人"
	for _, value := range str1 {
		fmt.Printf("value=%c\n", value)
	}
	for i := 0; i < len([]rune(str1)); i++ {
		fmt.Printf("value=%c\n", []rune(str1)[i])
	}
	/*字符串的拼接使用*/
	builder := strings.Builder{}
	builder.WriteString("hello")
	builder.WriteString("world")
	fmt.Println(builder.String())
}
