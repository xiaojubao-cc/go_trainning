package main

import (
	"fmt"
	"sync"
)

func incrementor() chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()
	return ch
}
func puller(ch chan int) chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		sun := 0
		for num := range ch {
			sun += num
		}
		out <- sun
	}()
	return out
}

func gen(num ...int) chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 0; i < len(num); i++ {
			ch <- num[i]
		}
	}()
	return ch
}

func multiply(ch chan int) chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for num := range ch {
			out <- num * num
		}
	}()
	return out
}

func merge(chList ...chan int) chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	wg.Add(len(chList))
	for _, ch := range chList {
		go func(ch chan int) {
			defer wg.Done()
			defer close(out)
			for num := range ch {
				out <- num
			}
		}(ch)
	}
	go func() {
		wg.Wait()
	}()
	return out
}
func main() {
	ch := puller(incrementor())
	for sum := range ch {
		fmt.Println(sum)
	}

	out := merge(multiply(gen(1, 2, 3, 4, 5)), multiply(gen(6, 7, 8, 9, 10)))
	for mul := range out {
		fmt.Println(mul)
	}
}
