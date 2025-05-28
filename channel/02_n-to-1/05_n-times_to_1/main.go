package main

import "fmt"

func main() {
	//创建多个生产者
	n := 10
	ch := make(chan int)
	done := make(chan bool)

	for i := 0; i < n; i++ {
		go func() {
			for i := 0; i < 10; i++ {
				ch <- i
			}
			done <- true
		}()
	}

	go func() {
		for i := 0; i < n; i++ {
			fmt.Println(<-done)
		}
		close(ch)
	}()

	for i := range ch {
		fmt.Println(i)
	}
}
