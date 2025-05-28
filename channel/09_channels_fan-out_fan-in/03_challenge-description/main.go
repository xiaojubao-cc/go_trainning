package main

import (
	"fmt"
	"sync"
)

func main() {
	for ch := range merge(fanOut(gen(), 10)...) {
		fmt.Println(ch)
	}
}

func gen() <-chan int {
	out := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			for j := 3; j < 13; j++ {
				out <- j
			}
		}
		close(out)
	}()
	return out
}
func fanOut(ch <-chan int, n int) []<-chan int {
	//在这里使用make时需要注意容量设置，初始的len需要配置为0，以便后续的append扩容
	//var outs []<-chan int
	outs := make([]<-chan int, 0, n)
	for i := 0; i < n; i++ {
		outs = append(outs, factorial(ch))
	}
	return outs
}

func factorial(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- fact(n)
		}
		close(out)
	}()
	return out
}

func fact(n int) int {
	total := 1
	for i := n; i > 0; i-- {
		total *= i
	}
	return total
}

func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
