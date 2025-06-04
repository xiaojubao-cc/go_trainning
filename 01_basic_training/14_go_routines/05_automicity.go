package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

/*原子类*/
func atomicitySum(count *int32, wg *sync.WaitGroup) {
	defer wg.Done()
	atomic.AddInt32(count, 1)
}
func main() {
	var wg sync.WaitGroup
	var count int32
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go atomicitySum(&count, &wg)
	}
	wg.Wait()
	fmt.Println(atomic.LoadInt32(&count))
}
