package main

import (
	"fmt"
	"sync"
)

func main() {
	for ch := range merge(sq(gen(2, 3)), sq(gen(2, 3))) {
		fmt.Println(ch)
	}
}

func gen(nums ...int) chan int {
	fmt.Printf("TYPE OF NUMS %T\n", nums)

	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func sq(in chan int) chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func merge(cs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	wg.Add(len(cs))
	outPut := func(ch <-chan int) {
		for c := range ch {
			out <- c
		}
		wg.Done()
	}

	for _, ch := range cs {
		go outPut(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
