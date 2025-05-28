package main

import "fmt"

func main() {
	ch := incrementor()
	sum := puller(ch)
	for t := range sum {
		fmt.Println(t)
	}
}

func incrementor() chan int {
	out := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			out <- i
		}
		close(out)
	}()
	return out
}

func puller(ch chan int) chan int {
	out := make(chan int)
	go func() {
		sum := 0
		for i := range ch {
			sum += i
		}
		out <- sum
		close(out)
	}()
	return out
}
