package main

import (
	"fmt"
	"sync"
)

/*
不能重复关闭通道，由发送方关闭通道;无缓冲通道需要发送和接收在不同的协程中准备就绪;避免协程间的相互等待
通道关闭后不能再接收数据，否则会panic;使用select注意添加default分支
*/
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
			for num := range ch {
				out <- num
			}
		}(ch)
	}
	/*
		非阻塞返回：merge 函数立即返回 out 通道，调用者可以立刻开始消费数据。
		通道关闭时机正确：在所有数据发送完毕后关闭通道，避免数据丢失。
		避免死锁：若 wg.Wait() 在主协程，merge 函数会阻塞，导致调用者无法接收数据，形成死锁。
	*/
	go func() {
		wg.Wait()
		close(out)
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
