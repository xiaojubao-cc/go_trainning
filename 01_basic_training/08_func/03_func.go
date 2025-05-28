package main

import "fmt"

/*可变参数必须放在最后*/
func multiplyByOdd(callback func(int) bool, is ...int) int {
	multi := 1
	for _, val := range is {
		if callback(val) {
			multi *= val
		}
	}
	return multi
}
func main() {
	odd := func(i int) bool {
		return i%2 == 1
	}
	var slice = []int{1, 2, 3, 4, 5}
	fmt.Printf("函数作为参数的函数返回值是： %d", multiplyByOdd(odd, slice...))
}
