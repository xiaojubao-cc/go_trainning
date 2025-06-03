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
			fmt.Println("任务取消")
			break LOOP
		default:
			fmt.Println("任务执行中")
		}
	}
	wg.Done()
}

func work1(cx context.Context) {
LOOP:
	for {
		fmt.Println("worker1")
		time.Sleep(time.Second)
		select {
		case <-cx.Done():
			//阻塞等待取消函数的调用通知
			fmt.Println("任务1取消")
			break LOOP
		default:
			fmt.Println("任务1执行中")
		}
	}
}

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	wg.Add(1)
	go work(ctx)
	fmt.Println("任务开始取消调用...")
	defer cancelFunc()
	wg.Wait()
	fmt.Println("over")
}
