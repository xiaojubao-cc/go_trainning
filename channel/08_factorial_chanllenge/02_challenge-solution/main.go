package main

import "fmt"

func main() {
	for i := range factorial(gen()) {
		fmt.Println(i)
	}
}

func gen() <-chan int {
	ch := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			for j := 3; j < 13; j++ {
				ch <- j
			}
		}
		close(ch)
	}()
	return ch
}

func factorial(ch <-chan int) <-chan int {
	read := make(chan int)
	go func() {
		for i := range ch {
			read <- fact(i)
		}
		close(read)
	}()
	return read
}

func fact(n int) int {
	total := 1
	for i := n; i > 0; i-- {
		total *= i
	}
	return total
}
