package main

/*withValue*/
import (
	"context"
	"fmt"
	"sync"
	"time"
)

type TraceCode string

var wg1 sync.WaitGroup

func worker(ctx context.Context) {
	/*主要业务*/
	key := TraceCode("TRACE_CODE")
	val, ok := ctx.Value(key).(string)
	if !ok {
		fmt.Errorf("invalid trace code")
	}
LOOP:
	for {
		fmt.Printf("worker, trace code:%s\n", val)
		time.Sleep(time.Millisecond * 10) // 假设正常连接数据库耗时10毫秒
		select {
		case <-ctx.Done(): // 50毫秒后自动调用
			break LOOP
		default:
		}
	}
	fmt.Println("worker done!")
	wg1.Done()
}
func main() {
	// 设置一个50毫秒的超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
	ctx = context.WithValue(ctx, TraceCode("TRACE_CODE"), "123456")
	wg1.Add(1)
	defer cancel()
	go worker(ctx)
	wg1.Wait()
}
