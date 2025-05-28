package main

import "fmt"

/*递归函数需要有退出递归的标识*/
func factorial(i int) int {
	if i == 0 {
		return 1
	}
	return i * factorial(i-1)

}
func main() {
	fmt.Printf("循环递归函数返回值：%d", factorial(5))
}
