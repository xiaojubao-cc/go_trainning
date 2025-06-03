package main

import (
	"fmt"
	"sync"
	"time"
)

/*信号量*/
func main() {
	maxConcurrent := 3
	sem := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			/*获取信号量*/
			sem <- struct{}{}
			defer func() {
				/*释放信号量*/
				<-sem
			}()
			/*执行操作*/
			fmt.Printf("任务 %d 开始\n", id)
			time.Sleep(time.Second)
			fmt.Printf("任务 %d 完成\n", id)
		}(i)
	}
	wg.Wait()
}
