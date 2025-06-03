package main

import (
	"fmt"
	"sync"
)

/*锁*/
func mutexMethod(count *int, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	mu.Lock()
	*count++
	mu.Unlock()
}

func main() {
	var mu sync.Mutex
	var wg sync.WaitGroup
	count := 0
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go mutexMethod(&count, &mu, &wg)
	}
	wg.Wait()
	fmt.Println(count)
}
