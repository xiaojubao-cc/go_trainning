package main

import "fmt"

func main() {
	//使用range进行取值

	ch := make(chan int)
	//启动协程
	go func() {
		//这里需要关闭channel不然会出现死锁
		defer close(ch)
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()

	for i := range ch {
		fmt.Println(i)
	}
}
