package main

import (
	"context"
	"fmt"
	"time"
)

/*
特性               WithDeadline                   WithTimeout
参数类型         绝对时间点（time.Time）        相对时间（time.Duration）
适用场景         需要固定时间点触发取消           需要动态设置超时时间
实现方式       内部调用 WithTimeout 实现      内部调用 WithDeadline 实现
*/
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
