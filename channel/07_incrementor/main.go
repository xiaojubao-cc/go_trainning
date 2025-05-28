package main

import "fmt"

func main() {
	//多个叠加器
	g1 := incrementor("goroutine01")
	g2 := incrementor("goroutine02")
	p1 := puller(g1)
	p2 := puller(g2)
	fmt.Println(<-p1)
	fmt.Println(<-p2)
}

func incrementor(s string) <-chan int {
	read := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			read <- i
			fmt.Println(s + string(i))
		}
		close(read)
	}()
	return read
}

func puller(ch <-chan int) <-chan int {
	read := make(chan int)
	go func() {
		sum := 0
		for i := range ch {
			sum += i
		}
		read <- sum
		close(read)
	}()
	return read
}
