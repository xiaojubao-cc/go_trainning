package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

/*原子类*/
func atomicitySum(count *int64, wg *sync.WaitGroup) {
	defer wg.Done()
	atomic.AddInt64(count, 1)
}
func main() {
	var wg sync.WaitGroup
	var count int64
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go atomicitySum(&count, &wg)
	}
	wg.Wait()
	fmt.Println(atomic.LoadInt64(&count))
}
