package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)
	//启动协程
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()

	//启动协程消费
	go func() {
		for {
			fmt.Println(<-ch)
		}
	}()

	//主协程进行阻塞
	time.Sleep(2 * time.Second)
}
