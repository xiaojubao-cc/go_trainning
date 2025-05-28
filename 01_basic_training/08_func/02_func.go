package main

import "fmt"

func add(a, b int) (s int) {
	defer func() {
		s -= 10
	}()
	s = a + b
	return
}

func multiply(is ...int) int {
	mult := 1
	for _, val := range is {
		mult *= val
	}
	return mult
}

func main() {
	fmt.Println(add(1, 2))
	var arr = []int{1, 2, 3, 4, 5}
	val := multiply(arr...)
	fmt.Printf("可变参数乘积：%d", val)
}
