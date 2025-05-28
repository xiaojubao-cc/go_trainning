package main

import "fmt"

func main() {
	ch := incrementor()
	for i := range puller(ch) {
		fmt.Println(i)
	}
}

func incrementor() <-chan int {
	//返回一个只读通道
	read := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			read <- i
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
