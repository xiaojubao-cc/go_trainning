package main

import (
	"fmt"
	"sync"
)

// Channel /*go并发管道适合协程间通信*/
type Channel chan interface{}

func main() {
	ch := make(Channel)
	go func() {
		/*关于管道关闭的时机，应该尽量在向管道发送数据的那一方关闭管道，而不要在接收方关闭管道，因为大多数情况下接收方只知道接收数据，并不知道该在什么时候关闭管道。*/
		defer close(ch)
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()
	for i := range ch {
		fmt.Println(i)
	}

	ch1 := make(chan int, 5)
	go func() {
		/*关于管道关闭的时机，应该尽量在向管道发送数据的那一方关闭管道，而不要在接收方关闭管道，因为大多数情况下接收方只知道接收数据，并不知道该在什么时候关闭管道。*/
		defer close(ch1)
		for i := 0; i < 10; i++ {
			ch1 <- i
		}
	}()
	for i := range ch1 {
		fmt.Println(i)
	}
	/*一个生产者，多个消费者*/
	ch2 := make(Channel)
	go func() {
		for i := 0; i < 100000; i++ {
			ch2 <- i
		}
		defer close(ch2)
	}()
	/*多个消费者*/
	var wg sync.WaitGroup
	sem := make(chan struct{}, 5)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			/*获取信号量*/
			sem <- struct{}{}
			defer func() {
				/*释放信号量*/
				<-sem
			}()
			/*这里进行数据消费*/
			for val := range ch2 {
				fmt.Printf("执行数据消费: %d\n", val)
			}
		}()
	}
	wg.Wait()

}
