package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type ResponseData struct {
	response *http.Response
	err      error
}

func call(ctx context.Context) {
	transport := http.Transport{
		// 请求频繁可定义全局的client对象并启用长链接
		// 请求不频繁使用短链接
		DisableKeepAlives: true}
	client := http.Client{
		Transport: &transport,
	}
	responseData := make(chan *ResponseData, 1)
	request, err := http.NewRequest("GET", "http://127.0.0.1:8080/", nil)
	if err != nil {
		fmt.Printf("new requestg failed, err:%v\n", err)
		return
	}
	request = request.WithContext(ctx)
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()
	go func() {
		response, err2 := client.Do(request)
		returnData := &ResponseData{response, err2}
		responseData <- returnData
		wg.Done()
	}()
	select {
	case <-ctx.Done():
		fmt.Println("timeout")
		return
	case result := <-responseData:
		fmt.Println("call server api success")
		if result.err != nil {
			fmt.Printf("call server api failed, err:%v\n", result.err)
			return
		}
		defer result.response.Body.Close()
		data, _ := ioutil.ReadAll(result.response.Body)
		fmt.Printf("resp:%v\n", string(data))
	}
}

func main() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancelFunc()
	call(ctx)
}
