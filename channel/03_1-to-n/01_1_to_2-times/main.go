package main

import (
	"fmt"
	"sync"
)

func main() {
	//ch := make(chan int)
	//done := make(chan bool)
	//
	//go func() {
	//	for i := 0; i < 100000; i++ {
	//		ch <- i
	//	}
	//	close(ch)
	//}()
	//
	//for i := 0; i < 10; i++ {
	//	go func() {
	//		for i := range ch {
	//			fmt.Println(i)
	//		}
	//		done <- true
	//	}()
	//}
	//
	//for i := 0; i < 10; i++ {
	//	<-done
	//}

	ch := make(chan int)
	var wg sync.WaitGroup

	// 生产者协程
	go func() {
		for i := 0; i < 100000; i++ {
			ch <- i
		}
		close(ch) // 关闭通道告知消费者数据发送完成
	}()

	// 启动 10 个消费者协程
	for i := 0; i < 10; i++ {
		wg.Add(1) // 每个协程启动前 +1
		go func() {
			defer wg.Done() // 确保协程退出时 -1（即使发生 panic）
			for num := range ch {
				fmt.Println(num)
			}
		}()
	}

	wg.Wait() // 主协程阻塞，直到所有消费者完成
}
