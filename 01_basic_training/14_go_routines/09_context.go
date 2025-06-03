package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

/*模拟客户端和服务端*/
func handler(resp http.ResponseWriter, req *http.Request) {
	number := rand.Intn(2)
	if number == 0 {
		time.Sleep(time.Second * 10) // 耗时10秒的慢响应
		fmt.Fprintf(resp, "slow response")
		return
	}
	fmt.Fprint(resp, "quick response")
}
func main() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Errorf("服务器异常：%s", err)
	}
}
