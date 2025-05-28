package main

import "fmt"

func main() {
	//使用另外一个信道来控制
	ch := make(chan int)
	done := make(chan bool)

	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		done <- true
	}()

	go func() {
		//创建了几个goroutine就会创建几个done
		<-done
		<-done
		close(ch)
	}()

	for i := range ch {
		fmt.Println(i)
	}

}
