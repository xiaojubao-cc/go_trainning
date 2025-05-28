package main

import "fmt"

/*
函数go中不支持函数的重载
函数的参数都是值传递，传递的都是副本信息
*/
func updateSlice(s []int) []int {
	s[0] = 100       //修改共享的底层数据会影响原值
	s = append(s, 6) //触发扩容操作，不再共享原数组，不影响原值
	return s
}

func updateStr(s string) {
	s = "hello"
}

/*闭包*/

func closures() func() int {
	var i int = 0
	return func() int {
		i++
		return i
	}
}
func main() {
	slice := make([]int, 0)
	slice = append(slice, 1, 2, 3, 4, 5)
	updateSlice(slice)
	for _, value := range slice {
		fmt.Println(value)
	}
	s := "golang"
	updateStr(s)
	fmt.Println(s)
	/*闭包*/
	returnFunc := closures()
	for i := 0; i < 10; i++ {
		fmt.Println(returnFunc())
	}
}
