package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	after := time.Now().Add(5 * time.Second)
	/*WithTimeout返回WithDeadline(parent, time.Now().Add(timeout))。*/
	deadline, cancelFunc := context.WithDeadline(context.Background(), after)
	defer cancelFunc()
	select {
	case <-time.After(1 * time.Second):
		fmt.Println("overslept")
	case <-deadline.Done():
		fmt.Println(deadline.Err())
	}
}
