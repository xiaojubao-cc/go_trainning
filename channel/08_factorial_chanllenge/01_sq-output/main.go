package main

import "fmt"

func main() {
	ch := gen(2, 3)
	for i := range factorial(factorial(ch)) {
		fmt.Println(i)
	}
	/*f1 := factorial(ch)
	fmt.Println(<-f1)
	fmt.Println(<-f1)*/
}

func gen(nums ...int) chan int {
	ch := make(chan int)
	go func() {
		for _, i := range nums {
			ch <- i
		}
		close(ch)
	}()
	return ch
}

func factorial(ch chan int) chan int {
	read := make(chan int)
	go func() {
		for i := range ch {
			read <- i * i
		}
		close(read)
	}()
	return read
}
