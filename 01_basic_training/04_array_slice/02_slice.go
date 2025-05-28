package main

import "fmt"

/*切片*/
func main() {
	var slice = make([]int, 0)
	for index, value := range append(slice, 1, 2, 3, 4, 5) {
		fmt.Println(index, value)
	}
	/*此处是操作指针*/
	var slice1 *[]int = new([]int)
	for index, value := range append(*slice1, 6, 7, 8, 9, 10) {
		fmt.Println(index, value)
	}

}
