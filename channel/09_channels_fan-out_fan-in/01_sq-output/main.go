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

func merge(cs ...chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	wg.Add(len(cs))
	for _, c := range cs {
		//这里使用闭包传值为了避免循环变量在并发中的共享问题
		go func(ch chan int) {
			for i := range ch {
				out <- i
			}
			wg.Done()
		}(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
