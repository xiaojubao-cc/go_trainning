package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

/*上下文父进程通知子进程*/
var wg sync.WaitGroup

func work(cx context.Context) {
	go work1(cx)
LOOP:
	for {
		fmt.Println("worker")
		time.Sleep(time.Second)
		select {
		case <-cx.Done():
			//阻塞等待取消函数的调用通知
			fmt.Println("work is canceled")
			break LOOP
		default:
			fmt.Println("work is running")
		}
	}
	wg.Done()
}

func work1(cx context.Context) {
	work2(cx)
LOOP:
	for {
		fmt.Println("worker1")
		time.Sleep(time.Second)
		select {
		case <-cx.Done():
			//阻塞等待取消函数的调用通知
			fmt.Println("work1 is canceled")
			break LOOP
		default:
			fmt.Println("work1 is running")
		}
	}
}
func work2(ctx context.Context) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("work2 is canceled")
			break LOOP
		default:
			fmt.Println("work2 is running")
		}
	}
}
func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	wg.Add(1)
	go work(ctx)
	fmt.Println("task is canceling...")
	cancelFunc()
	wg.Wait()
	fmt.Println("task is over")
}
