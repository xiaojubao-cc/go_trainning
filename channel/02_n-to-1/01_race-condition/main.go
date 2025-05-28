package main

import (
	"fmt"
	"sync"
)

func main() {
	//多个生产的goroutine,单个消费的goroutine

	ch := make(chan int)
	//启动多个生产的goroutine
	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		for i := 0; i < 10; i++ {
			ch <- i
		}
		wg.Done()
	}()

	go func() {
		wg.Add(1)
		for i := 0; i < 10; i++ {
			ch <- i
		}
		wg.Done()
	}()

	//启动协程进行关闭组避免出现死锁
	go func() {
		wg.Wait()
		close(ch)
	}()

	for i := range ch {
		fmt.Println(i)
	}
}
